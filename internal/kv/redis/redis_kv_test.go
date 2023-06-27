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

package rediskv_test

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/cockroachdb/errors"
	rdsclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	rediskv "github.com/milvus-io/milvus/internal/kv/redis"
	"github.com/milvus-io/milvus/pkg/util/funcutil"
	"github.com/milvus-io/milvus/pkg/util/paramtable"
)

var Params paramtable.ComponentParam
var c *rdsclient.Client

func setupRedis() {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	c = rdsclient.NewClient(&rdsclient.Options{
		Addr: s.Addr(),
	})
}

func setupRemoteRedis() {
	pwd := os.Getenv("RedisPwd")
	if len(pwd) == 0 {
		panic("Please set redis password to environment variable RedisPwd")
	}
	addr := os.Getenv("RedisAddr")
	if len(addr) == 0 {
		panic("Please set redis address with port to environment variable RedisAddr")
	}
	c = rdsclient.NewClient(&rdsclient.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0, // use default DB
	})
}

func TestMain(m *testing.M) {
	Params.Init()
	setupRedis()
	code := m.Run()
	os.Exit(code)
}

func TestRedisKV_Load(te *testing.T) {
	te.Run("kv SaveAndLoad", func(t *testing.T) {
		rootPath := "/redis/test/root/saveandload"
		kv := rediskv.NewRedisKV(c, rootPath)
		err := kv.RemoveWithPrefix("")
		require.NoError(t, err)

		defer kv.Close()
		defer kv.RemoveWithPrefix("")

		saveAndLoadTests := []struct {
			key   string
			value string
		}{
			{"test1", "value1"},
			{"test2", "value2"},
			{"test1/a", "value_a"},
			{"test1/b", "value_b"},
		}

		for i, test := range saveAndLoadTests {
			if i < 4 {
				err = kv.Save(test.key, test.value)
				assert.NoError(t, err)
			}

			val, err := kv.Load(test.key)
			assert.NoError(t, err)
			assert.Equal(t, test.value, val)
		}

		invalidLoadTests := []struct {
			invalidKey string
		}{
			{"t"},
			{"a"},
			{"test1a"},
		}

		for _, test := range invalidLoadTests {
			val, err := kv.Load(test.invalidKey)
			assert.Error(t, err)
			assert.Zero(t, val)
		}

		loadPrefixTests := []struct {
			prefix string

			expectedKeys   []string
			expectedValues []string
			expectedError  error
		}{
			{"test", []string{
				kv.GetPath("test1"),
				kv.GetPath("test2"),
				kv.GetPath("test1/a"),
				kv.GetPath("test1/b")}, []string{"value1", "value2", "value_a", "value_b"}, nil},
			{"test1", []string{
				kv.GetPath("test1"),
				kv.GetPath("test1/a"),
				kv.GetPath("test1/b")}, []string{"value1", "value_a", "value_b"}, nil},
			{"test2", []string{kv.GetPath("test2")}, []string{"value2"}, nil},
			{"", []string{
				kv.GetPath("test1"),
				kv.GetPath("test2"),
				kv.GetPath("test1/a"),
				kv.GetPath("test1/b")}, []string{"value1", "value2", "value_a", "value_b"}, nil},
			{"test1/a", []string{kv.GetPath("test1/a")}, []string{"value_a"}, nil},
			{"a", []string{}, []string{}, nil},
			{"root", []string{}, []string{}, nil},
			{"/redis/test/root", []string{}, []string{}, nil},
		}

		for _, test := range loadPrefixTests {
			actualKeys, actualValues, err := kv.LoadWithPrefix(test.prefix)
			assert.ElementsMatch(t, test.expectedKeys, actualKeys)
			assert.ElementsMatch(t, test.expectedValues, actualValues)
			assert.Equal(t, test.expectedError, err)
		}

		removeTests := []struct {
			validKey   string
			invalidKey string
		}{
			{"test1", "abc"},
			{"test1/a", "test1/lskfjal"},
			{"test1/b", "test1/b"},
			{"test2", "-"},
		}

		for _, test := range removeTests {
			err = kv.Remove(test.validKey)
			assert.NoError(t, err)

			_, err = kv.Load(test.validKey)
			assert.Error(t, err)

			err = kv.Remove(test.validKey)
			assert.NoError(t, err)
			err = kv.Remove(test.invalidKey)
			assert.NoError(t, err)
		}
	})

	te.Run("kv MultiSaveAndMultiLoad", func(t *testing.T) {
		rootPath := "/redis/test/root/multi_save_and_multi_load"
		kv := rediskv.NewRedisKV(c, rootPath)

		defer kv.Close()
		defer kv.RemoveWithPrefix("")

		multiSaveTests := map[string]string{
			"key_1":      "value_1",
			"key_2":      "value_2",
			"key_3/a":    "value_3a",
			"multikey_1": "multivalue_1",
			"multikey_2": "multivalue_2",
			"_":          "other",
		}

		err := kv.MultiSave(multiSaveTests)
		assert.NoError(t, err)
		for k, v := range multiSaveTests {
			actualV, err := kv.Load(k)
			assert.NoError(t, err)
			assert.Equal(t, v, actualV)
		}

		multiLoadTests := []struct {
			inputKeys      []string
			expectedValues []string
		}{
			{[]string{"key_1"}, []string{"value_1"}},
			{[]string{"key_1", "key_2", "key_3/a"}, []string{"value_1", "value_2", "value_3a"}},
			{[]string{"multikey_1", "multikey_2"}, []string{"multivalue_1", "multivalue_2"}},
			{[]string{"_"}, []string{"other"}},
		}

		for _, test := range multiLoadTests {
			vs, err := kv.MultiLoad(test.inputKeys)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValues, vs)
		}

		invalidMultiLoad := []struct {
			invalidKeys    []string
			expectedValues []string
		}{
			{[]string{"a", "key_1"}, []string{"", "value_1"}},
			{[]string{".....", "key_1"}, []string{"", "value_1"}},
			{[]string{"*********"}, []string{""}},
			{[]string{"key_1", "1"}, []string{"value_1", ""}},
		}

		for _, test := range invalidMultiLoad {
			vs, err := kv.MultiLoad(test.invalidKeys)
			assert.Error(t, err)
			assert.Equal(t, test.expectedValues, vs)
		}

		removeWithPrefixTests := []string{
			"key_1",
			"multi",
		}

		for _, k := range removeWithPrefixTests {
			err = kv.RemoveWithPrefix(k)
			assert.NoError(t, err)

			ks, vs, err := kv.LoadWithPrefix(k)
			assert.Empty(t, ks)
			assert.Empty(t, vs)
			assert.NoError(t, err)
		}

		multiRemoveTests := []string{
			"key_2",
			"key_3/a",
			"multikey_2",
			"_",
		}

		err = kv.MultiRemove(multiRemoveTests)
		assert.NoError(t, err)

		ks, vs, err := kv.LoadWithPrefix("")
		assert.NoError(t, err)
		assert.Empty(t, ks)
		assert.Empty(t, vs)

		multiSaveAndRemoveTests := []struct {
			multiSaves   map[string]string
			multiRemoves []string
		}{
			{map[string]string{"key_1": "value_1"}, []string{}},
			{map[string]string{"key_2": "value_2"}, []string{"key_1"}},
			{map[string]string{"key_3/a": "value_3a"}, []string{"key_2"}},
			{map[string]string{"multikey_1": "multivalue_1"}, []string{}},
			{map[string]string{"multikey_2": "multivalue_2"}, []string{"multikey_1", "key_3/a"}},
			{make(map[string]string), []string{"multikey_2"}},
		}
		for _, test := range multiSaveAndRemoveTests {
			err = kv.MultiSaveAndRemove(test.multiSaves, test.multiRemoves)
			assert.NoError(t, err)
		}

		ks, vs, err = kv.LoadWithPrefix("")
		assert.NoError(t, err)
		assert.Empty(t, ks)
		assert.Empty(t, vs)
	})

	te.Run("kv MultiRemoveWithPrefix", func(t *testing.T) {
		rootPath := "/redis/test/root/multi_remove_with_prefix"
		kv := rediskv.NewRedisKV(c, rootPath)
		defer kv.Close()
		defer kv.RemoveWithPrefix("")

		prepareTests := map[string]string{
			"x/abc/1": "1",
			"x/abc/2": "2",
			"x/def/1": "10",
			"x/def/2": "20",
			"x/den/1": "100",
			"x/den/2": "200",
		}

		err := kv.MultiSave(prepareTests)
		require.NoError(t, err)

		multiRemoveWithPrefixTests := []struct {
			prefix []string

			testKey       string
			expectedValue string
		}{
			{[]string{"x/abc"}, "x/abc/1", ""},
			{[]string{}, "x/abc/2", ""},
			{[]string{}, "x/def/1", "10"},
			{[]string{}, "x/def/2", "20"},
			{[]string{}, "x/den/1", "100"},
			{[]string{}, "x/den/2", "200"},
			{[]string{}, "not-exist", ""},
			{[]string{"x/def", "x/den"}, "x/def/1", ""},
			{[]string{}, "x/def/1", ""},
			{[]string{}, "x/def/2", ""},
			{[]string{}, "x/den/1", ""},
			{[]string{}, "x/den/2", ""},
			{[]string{}, "not-exist", ""},
		}

		for _, test := range multiRemoveWithPrefixTests {
			if len(test.prefix) > 0 {
				err = kv.MultiRemoveWithPrefix(test.prefix)
				assert.NoError(t, err)
			}

			v, _ := kv.Load(test.testKey)
			assert.Equal(t, test.expectedValue, v)
		}

		k, v, err := kv.LoadWithPrefix("/")
		assert.NoError(t, err)
		assert.Zero(t, len(k))
		assert.Zero(t, len(v))

		// MultiSaveAndRemoveWithPrefix
		err = kv.MultiSave(prepareTests)
		require.NoError(t, err)
		multiSaveAndRemoveWithPrefixTests := []struct {
			multiSave map[string]string
			prefix    []string

			loadPrefix         string
			lengthBeforeRemove int
			lengthAfterRemove  int
		}{
			{map[string]string{}, []string{"x/abc", "x/def", "x/den"}, "x", 6, 0},
			{map[string]string{"y/a": "vvv", "y/b": "vvv"}, []string{}, "y", 0, 2},
			{map[string]string{"y/c": "vvv"}, []string{}, "y", 2, 3},
			{map[string]string{"p/a": "vvv"}, []string{"y/a", "y"}, "y", 3, 0},
			{map[string]string{}, []string{"p"}, "p", 1, 0},
		}

		for _, test := range multiSaveAndRemoveWithPrefixTests {
			k, _, err = kv.LoadWithPrefix(test.loadPrefix)
			assert.NoError(t, err)
			assert.Equal(t, test.lengthBeforeRemove, len(k))

			err = kv.MultiSaveAndRemoveWithPrefix(test.multiSave, test.prefix)
			assert.NoError(t, err)

			k, _, err = kv.LoadWithPrefix(test.loadPrefix)
			assert.NoError(t, err)
			assert.Equal(t, test.lengthAfterRemove, len(k))
		}
	})
}

func Test_WalkWithPagination(t *testing.T) {
	rootPath := "/redis/test/root/pagination"
	kv := rediskv.NewRedisKV(c, rootPath)

	defer kv.Close()
	defer kv.RemoveWithPrefix("")

	kvs := map[string]string{
		"A/100":    "v1",
		"AA/100":   "v2",
		"AB/100":   "v3",
		"AB/2/100": "v4",
		"B/100":    "v5",
	}

	err := kv.MultiSave(kvs)
	assert.NoError(t, err)
	for k, v := range kvs {
		actualV, err := kv.Load(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actualV)
	}

	t.Run("apply function error ", func(t *testing.T) {
		err = kv.WalkWithPrefix("A", 5, func(key []byte, value []byte) error {
			return errors.New("error")
		})
		assert.Error(t, err)
	})

	t.Run("get with non-exist prefix ", func(t *testing.T) {
		err = kv.WalkWithPrefix("non-exist-prefix", 5, func(key []byte, value []byte) error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("with different pagination", func(t *testing.T) {
		testFn := func(pagination int) {
			expected := map[string]string{
				"A/100":    "v1",
				"AA/100":   "v2",
				"AB/100":   "v3",
				"AB/2/100": "v4",
			}

			expectedKeys := maps.Keys(expected)
			sort.Strings(expectedKeys)

			ret := make(map[string]string)
			actualKeys := make([]string, 0)

			err = kv.WalkWithPrefix("A", pagination, func(key []byte, value []byte) error {
				k := string(key)
				k = k[len(rootPath)+1:]
				ret[k] = string(value)
				actualKeys = append(actualKeys, k)
				return nil
			})

			assert.NoError(t, err)
			assert.Equal(t, expected, ret, fmt.Errorf("pagination: %d", pagination))
			// Ignore the order.
			assert.ElementsMatch(t, expectedKeys, actualKeys, fmt.Errorf("pagination: %d", pagination))
		}

		for p := -1; p < 6; p++ {
			testFn(p)
		}
		testFn(-100)
		testFn(100)
	})
}

func TestElapse(t *testing.T) {
	start := time.Now()
	isElapse := rediskv.CheckElapseAndWarn(start, "err message")
	assert.Equal(t, isElapse, false)

	time.Sleep(2001 * time.Millisecond)
	isElapse = rediskv.CheckElapseAndWarn(start, "err message")
	assert.Equal(t, isElapse, true)
}

func TestCheckValueSizeAndWarn(t *testing.T) {
	ret := rediskv.CheckValueSizeAndWarn("k", "v")
	assert.False(t, ret)

	v := make([]byte, 1024000)
	ret = rediskv.CheckValueSizeAndWarn("k", v)
	assert.True(t, ret)
}

func TestCheckTnxBytesValueSizeAndWarn(t *testing.T) {
	kvs := make(map[string][]byte, 0)
	kvs["k"] = []byte("v")
	ret := rediskv.CheckTnxBytesValueSizeAndWarn(kvs)
	assert.False(t, ret)

	kvs["k"] = make([]byte, 1024000)
	ret = rediskv.CheckTnxBytesValueSizeAndWarn(kvs)
	assert.True(t, ret)
}

func TestCheckTnxStringValueSizeAndWarn(t *testing.T) {
	kvs := make(map[string]string, 0)
	kvs["k"] = "v"
	ret := rediskv.CheckTnxStringValueSizeAndWarn(kvs)
	assert.False(t, ret)

	kvs["k1"] = funcutil.RandomString(1024000)
	ret = rediskv.CheckTnxStringValueSizeAndWarn(kvs)
	assert.True(t, ret)
}
