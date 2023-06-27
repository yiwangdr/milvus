// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rediskv

import (
	"context"
	"encoding/binary"
	"fmt"
	"path"
	"time"

	"github.com/cockroachdb/errors"
	rdsclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/milvus-io/milvus/internal/kv"
	"github.com/milvus-io/milvus/pkg/common"
	"github.com/milvus-io/milvus/pkg/log"
	/*
		"github.com/milvus-io/milvus/pkg/metrics"
		"github.com/milvus-io/milvus/pkg/util/timerecord"
	*/)

const (
	// RequestTimeout is default timeout for redis request.
	RequestTimeout = 10 * time.Second
)

// implementation assertion
var _ kv.MetaKv = (*redisKV)(nil)

// redisKV implements TxnKV interface, it supports to process multiple kvs in a transaction.
type redisKV struct {
	client   *rdsclient.Client
	rootPath string
}

// NewRedisKV creates a new redis kv.
func NewRedisKV(client *rdsclient.Client, rootPath string) *redisKV {
	kv := &redisKV{
		client:   client,
		rootPath: rootPath,
	}
	return kv
}

// Close closes the connection to redis.
func (kv *redisKV) Close() {
	log.Debug("redis kv closed", zap.String("path", kv.rootPath))
}

// GetPath returns the path of the key.
func (kv *redisKV) GetPath(key string) string {
	return path.Join(kv.rootPath, key)
}

func (kv *redisKV) WalkWithPrefix(prefix string, paginationSize int, fn func([]byte, []byte) error) error {
	start := time.Now()
	prefix = path.Join(kv.rootPath, prefix)
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	cursor := uint64(0)
	keySize := paginationSize
	if keySize < 0 {
		keySize = 0
	}
	batchKeys := make([]string, 0, keySize)

	for {
		var err error
		// Note that the walk is not done in a sorted order.
		// TODO: (1) confirm if sort is required.
		//       (2) check if scan pagination breaks on newly inserted records.
		batchKeys, cursor, err = kv.client.Scan(ctx, cursor, prefix+"*", int64(paginationSize)).Result()
		if err != nil {
			return errors.Wrap(err, "Failed redis scan while walking")
		}

		for _, key := range batchKeys {
			value, err := kv.client.Get(ctx, key).Bytes()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to read key %s while walking", key))
			}

			if err := fn([]byte(key), value); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to apply func on kv(%s, %s) while walking", key, value))
			}
		}

		if cursor == 0 {
			break
		}
	}

	CheckElapseAndWarn(start, "Slow redis operation(WalkWithPagination)", zap.String("prefix", prefix))
	return nil
}

// Load returns value of the key.
func (kv *redisKV) Load(key string) (string, error) {
	start := time.Now()
	key = path.Join(kv.rootPath, key)
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()
	val, err := kv.client.Get(ctx, key).Result()
	if err != nil {
		if err == rdsclient.Nil {
			return "", common.NewKeyNotExistError(key)
		}
		return "", errors.Wrap(err, fmt.Sprintf("Failed to read key %s while walking", key))
	}
	CheckElapseAndWarn(start, "Slow redis operation load", zap.String("key", key))
	return val, nil
}

// MultiLoad gets the values of the keys in a transaction.
func (kv *redisKV) MultiLoad(keys []string) ([]string, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	// Start a transaction
	tx := kv.client.TxPipeline()

	// Add Get commands for each key to the transaction
	getCmds := make([]*rdsclient.StringCmd, len(keys))
	for i, key := range keys {
		getCmds[i] = tx.Get(ctx, path.Join(kv.rootPath, key))
	}

	// Execute the transaction. Don't early return on Exec error as we could get partial results.
	_, err := tx.Exec(ctx)

	// Extract the results from Get commands
	values := make([]string, len(keys))
	for i, cmd := range getCmds {
		value, err := cmd.Result()
		if err == rdsclient.Nil {
			values[i] = ""
		} else if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to get value for key %s in MultiLoad", keys[i]))
		} else {
			values[i] = value
		}
	}
	CheckElapseAndWarn(start, "Slow redis operation multi load", zap.Any("keys", keys))
	return values, err
}

// LoadWithPrefix returns all the keys and values with the given key prefix.
func (kv *redisKV) LoadWithPrefix(prefix string) ([]string, []string, error) {
	start := time.Now()
	prefix = path.Join(kv.rootPath, prefix)
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	// Get all keys with the given prefix
	keysCmd := kv.client.Keys(ctx, prefix+"*")
	if err := keysCmd.Err(); err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to generate cmd for keys with prefix %s for LoadWithPrefix", prefix))
	}

	keys := keysCmd.Val()
	if len(keys) == 0 {
		return nil, nil, nil
	}

	tx := kv.client.TxPipeline()
	getCmds := make([]*rdsclient.StringCmd, len(keys))
	for i, key := range keys {
		getCmds[i] = tx.Get(ctx, key)
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to exec LoadWithPrefix transaction")
	}

	// Extract the results from Get commands
	values := make([]string, len(keys))
	for i, cmd := range getCmds {
		value, err := cmd.Result()
		if err == rdsclient.Nil {
			values[i] = ""
		} else if err != nil {
			return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to get value for key %s in LoadWithPrefix", keys[i]))
		} else {
			values[i] = value
		}
	}
	CheckElapseAndWarn(start, "Slow redis operation load with prefix", zap.Strings("keys", keys))
	return keys, values, nil
}

// Save saves the key-value pair.
func (kv *redisKV) Save(key, value string) error {
	start := time.Now()
	key = path.Join(kv.rootPath, key)
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()
	CheckValueSizeAndWarn(key, value)
	err := kv.client.Set(ctx, key, value, RequestTimeout).Err()
	CheckElapseAndWarn(start, "Slow redis operation save", zap.String("key", key))
	return err
}

// MultiSave saves the key-value pairs in a transaction.
func (kv *redisKV) MultiSave(kvs map[string]string) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	tx := kv.client.TxPipeline()
	for key, value := range kvs {
		key = path.Join(kv.rootPath, key)
		tx.Set(ctx, key, value, 0)
	}

	_, err := tx.Exec(ctx)
	if err != nil {
		log.Warn("RedisKV MultiSave error", zap.Any("kvs", kvs), zap.Int("len", len(kvs)), zap.Error(err))
	}
	CheckElapseAndWarn(start, "Slow redis operation multi save", zap.Any("kvs", kvs))
	return err
}

// Remove removes the key.
func (kv *redisKV) Remove(key string) error {
	start := time.Now()
	key = path.Join(kv.rootPath, key)
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	err := kv.client.Del(ctx, key).Err()
	CheckElapseAndWarn(start, "Slow redis operation remove", zap.String("key", key))
	return err
}

// MultiRemove removes the keys in a transaction.
func (kv *redisKV) MultiRemove(keys []string) error {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	tx := kv.client.TxPipeline()
	for _, key := range keys {
		tx.Del(ctx, path.Join(kv.rootPath, key))
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		log.Warn("RedisKV MultiRemove error", zap.Strings("keys", keys), zap.Int("len", len(keys)), zap.Error(err))
	}
	CheckElapseAndWarn(start, "Slow redis operation multi remove", zap.Strings("keys", keys))
	return err
}

// RemoveWithPrefix removes the keys with given prefix.
func (kv *redisKV) RemoveWithPrefix(prefix string) error {
	start := time.Now()
	prefix = path.Join(kv.rootPath, prefix)
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	keysCmd := kv.client.Keys(ctx, prefix+"*")
	if err := keysCmd.Err(); err != nil {
		return err
	}

	keys := keysCmd.Val()
	if len(keys) == 0 {
		return nil
	}
	tx := kv.client.TxPipeline()
	for _, key := range keys {
		tx.Del(ctx, key)
	}
	_, err := tx.Exec(ctx)

	CheckElapseAndWarn(start, "Slow redis operation remove with prefix", zap.String("prefix", prefix))
	return err
}

// MultiSaveAndRemove saves the key-value pairs and removes the keys in a transaction.
func (kv *redisKV) MultiSaveAndRemove(saves map[string]string, removals []string) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	tx := kv.client.TxPipeline()
	// Add Set commands for each key-value pair to the transaction
	for key, value := range saves {
		tx.Set(ctx, path.Join(kv.rootPath, key), value, 0)
	}
	// Add Del commands for each key to be removed to the transaction
	for _, key := range removals {
		tx.Del(ctx, path.Join(kv.rootPath, key))
	}
	// Execute the transaction
	_, err := tx.Exec(ctx)
	if err != nil {
		log.Warn("RedisKV MultiSaveAndRemove error",
			zap.Any("saves", saves),
			zap.Strings("removes", removals),
			zap.Int("saveLength", len(saves)),
			zap.Int("removeLength", len(removals)),
			zap.Error(err))
	}
	CheckElapseAndWarn(start, "Slow redis operation multi save and remove", zap.Any("saves", saves), zap.Strings("removals", removals))
	return err
}

// MultiRemoveWithPrefix removes the keys with given prefix.
func (kv *redisKV) MultiRemoveWithPrefix(prefixes []string) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	tx := kv.client.TxPipeline()
	for _, prefix := range prefixes {
		prefix = path.Join(kv.rootPath, prefix)
		// Get all keys with the current prefix
		keysCmd := kv.client.Keys(ctx, prefix+"*")
		if err := keysCmd.Err(); err != nil {
			return err
		}

		prefixKeys := keysCmd.Val()

		// Add Del commands for each key with the current prefix to the transaction
		for _, key := range prefixKeys {
			tx.Del(ctx, key)
		}
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		log.Warn("RedisKV MultiRemoveWithPrefix error", zap.Strings("keys", prefixes), zap.Int("len", len(prefixes)), zap.Error(err))
	}
	CheckElapseAndWarn(start, "Slow redis operation multi remove with prefix", zap.Strings("keys", prefixes))
	return err
}

// MultiSaveAndRemoveWithPrefix saves kv in @saves and removes the keys with given prefix in @removals.
func (kv *redisKV) MultiSaveAndRemoveWithPrefix(saves map[string]string, removals []string) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.TODO(), RequestTimeout)
	defer cancel()

	tx := kv.client.TxPipeline()
	for key, value := range saves {
		tx.Set(ctx, path.Join(kv.rootPath, key), value, 0)
	}
	for _, prefix := range removals {
		prefix = path.Join(kv.rootPath, prefix)
		// Get all keys with the current prefix
		keysCmd := kv.client.Keys(ctx, prefix+"*")
		if err := keysCmd.Err(); err != nil {
			return err
		}
		prefixKeys := keysCmd.Val()
		// Add Del commands for each key with the current prefix to the transaction
		for _, key := range prefixKeys {
			tx.Del(ctx, key)
		}
	}

	// Execute the transaction
	_, err := tx.Exec(ctx)
	if err != nil {
		log.Warn("RedisKV MultiSaveAndRemoveWithPrefix error",
			zap.Any("saves", saves),
			zap.Strings("removes", removals),
			zap.Int("saveLength", len(saves)),
			zap.Int("removeLength", len(removals)),
			zap.Error(err))
	}
	CheckElapseAndWarn(start, "Slow redis operation multi save and move with prefix", zap.Any("saves", saves), zap.Strings("removals", removals))
	return err
}

// CheckElapseAndWarn checks the elapsed time and warns if it is too long.
func CheckElapseAndWarn(start time.Time, message string, fields ...zap.Field) bool {
	elapsed := time.Since(start)
	if elapsed.Milliseconds() > 2000 {
		log.Warn(message, append([]zap.Field{zap.String("time spent", elapsed.String())}, fields...)...)
		return true
	}
	return false
}

func CheckValueSizeAndWarn(key string, value interface{}) bool {
	size := binary.Size(value)
	if size > 102400 {
		log.Warn("value size large than 100kb", zap.String("key", key), zap.Int("value_size(kb)", size/1024))
		return true
	}
	return false
}

func CheckTnxBytesValueSizeAndWarn(kvs map[string][]byte) bool {
	var hasWarn bool
	for key, value := range kvs {
		if CheckValueSizeAndWarn(key, value) {
			hasWarn = true
		}
	}
	return hasWarn
}

func CheckTnxStringValueSizeAndWarn(kvs map[string]string) bool {
	newKvs := make(map[string][]byte, len(kvs))
	for key, value := range kvs {
		newKvs[key] = []byte(value)
	}

	return CheckTnxBytesValueSizeAndWarn(newKvs)
}

/*

func (kv *redisKV) removeEtcdMeta(ctx context.Context, key string, opts ...rdsclient.OpOption) (*rdsclient.DeleteResponse, error) {
	ctx1, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()

	start := timerecord.NewTimeRecorder("removeEtcdMeta")
	resp, err := kv.client.Delete(ctx1, key, opts...)
	elapsed := start.ElapseSpan()
	metrics.MetaOpCounter.WithLabelValues(metrics.MetaRemoveLabel, metrics.TotalLabel).Inc()

	if err == nil {
		metrics.MetaRequestLatency.WithLabelValues(metrics.MetaRemoveLabel).Observe(float64(elapsed.Milliseconds()))
		metrics.MetaOpCounter.WithLabelValues(metrics.MetaRemoveLabel, metrics.SuccessLabel).Inc()
	} else {
		metrics.MetaOpCounter.WithLabelValues(metrics.MetaRemoveLabel, metrics.FailLabel).Inc()
	}

	return resp, err
}

func (kv *redisKV) getTxnWithCmp(ctx context.Context, cmp ...rdsclient.Cmp) rdsclient.Txn {
	return kv.client.Txn(ctx).If(cmp...)
}

func (kv *redisKV) executeTxn(txn rdsclient.Txn, ops ...rdsclient.Op) (*rdsclient.TxnResponse, error) {
	start := timerecord.NewTimeRecorder("executeTxn")

	resp, err := txn.Then(ops...).Commit()
	elapsed := start.ElapseSpan()
	metrics.MetaOpCounter.WithLabelValues(metrics.MetaTxnLabel, metrics.TotalLabel).Inc()

	if err == nil && resp.Succeeded {
		// cal put meta kv size
		totalPutSize := 0
		for _, op := range ops {
			if op.IsPut() {
				totalPutSize += binary.Size(op.ValueBytes())
			}
		}
		metrics.MetaKvSize.WithLabelValues(metrics.MetaPutLabel).Observe(float64(totalPutSize))

		// cal get meta kv size
		totalGetSize := 0
		for _, rp := range resp.Responses {
			if rp.GetResponseRange() != nil {
				for _, v := range rp.GetResponseRange().Kvs {
					totalGetSize += binary.Size(v)
				}
			}
		}
		metrics.MetaKvSize.WithLabelValues(metrics.MetaGetLabel).Observe(float64(totalGetSize))
		metrics.MetaRequestLatency.WithLabelValues(metrics.MetaTxnLabel).Observe(float64(elapsed.Milliseconds()))
		metrics.MetaOpCounter.WithLabelValues(metrics.MetaTxnLabel, metrics.SuccessLabel).Inc()
	} else {
		metrics.MetaOpCounter.WithLabelValues(metrics.MetaTxnLabel, metrics.FailLabel).Inc()
	}

	return resp, err
}
*/
