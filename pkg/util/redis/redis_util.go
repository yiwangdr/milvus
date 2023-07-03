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

package redis

import (
	"os"
	"sync"

	"github.com/alicebob/miniredis"
	rdsclient "github.com/redis/go-redis/v9"
)

var (
	maxTxnNum = 128
	initOnce   sync.Once
	redisServer *miniredis.Miniredis
)


func setupEmbededRedis() *rdsclient.Client {
	initOnce.Do(func() {
		var err error
		redisServer, err = miniredis.Run()
		if err != nil {
			panic(err)
		}
	})
	if redisServer == nil {
		panic("MiniRedis server failed initialization!")
	}

	return rdsclient.NewClient(&rdsclient.Options{
		Addr: redisServer.Addr(),
	})
}

func setupRemoteRedis() *rdsclient.Client {
	pwd := os.Getenv("RedisPwd")
	if len(pwd) == 0 {
		panic("Please set redis password to environment variable RedisPwd")
	}
	addr := os.Getenv("RedisAddr")
	if len(addr) == 0 {
		panic("Please set redis address with port to environment variable RedisAddr")
	}
	return rdsclient.NewClient(&rdsclient.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0, // use default DB
	})
}

// GetRedisClient returns redis client
func GetRedisClient() (*rdsclient.Client, error) {
	return setupEmbededRedis(), nil
}
