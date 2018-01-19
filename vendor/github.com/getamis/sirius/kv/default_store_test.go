// +build go1.9

// Copyright 2017 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AtomicPut is not supported
//
// func TestDefaultStoreAtomicPut(t *testing.T) {
// 	kv := New()

// 	previous, err := kv.AtomicPut("test-key-1", []byte("test-value-1"))
// 	assert.NoError(t, err, "should be no error")
// 	assert.Nil(t, previous, "should be nil")

// 	previous, err = kv.AtomicPut("test-key-1", []byte("test-value-2"))
// 	assert.NoError(t, err, "should be no error")
// 	assert.Equal(t, previous.Value, []byte("test-value-1"), "should be equal")

// 	previous, err = kv.AtomicPut("test-key-1", []byte("test-value-3"))
// 	assert.NoError(t, err, "should be no error")
// 	assert.Equal(t, previous.Value, []byte("test-value-2"), "should be equal")
// }

func TestDefaultStorePutGet(t *testing.T) {
	kv := New()

	err := kv.Put("test-key-1", []byte("test-value-1"))
	assert.NoError(t, err, "should be no error")

	exist, err := kv.Exists("test-key-1")
	assert.True(t, exist, "should be true")
	assert.NoError(t, err, "should be no error")

	v, err := kv.Get("test-key-1")
	assert.NoError(t, err, "should be no error")
	assert.Equal(t, v.Value, []byte("test-value-1"), "should be equal")

	_, err = kv.List("test-key-1")
	assert.Equal(t, ErrNotSupported, err, "should be equal")
}

func TestDefaultStoreDelete(t *testing.T) {
	kv := New()

	testData := &KeyValue{
		Key:   "test-key-1",
		Value: []byte("test-value-1"),
	}

	err := kv.Put(testData.Key, testData.Value)
	assert.NoError(t, err, "should be no error")

	err = kv.Delete("test-key-1")
	assert.NoError(t, err, "should be no error")

	exist, err := kv.Exists("test-key-1")
	assert.False(t, exist, "should be false")
	assert.NoError(t, err, "should be no error")

	err = kv.Delete("test-key-1")
	assert.Equal(t, ErrKeyNotFound, err, "should be equal")

	// AtomicDelete is not supported
	//
	// testData = &KeyValue{
	// 	Key:   "test-key-2",
	// 	Value: []byte("test-value-2"),
	// }

	// err = kv.Put("test-key-2", []byte("test-value-2"))
	// assert.NoError(t, err, "should be no error")

	// v, err := kv.AtomicDelete("test-key-2")
	// assert.Equal(t, testData.Key, v.Key, "should be equal")
	// assert.Equal(t, testData.Value, v.Value, "should be equal")
	// assert.NoError(t, err, "should be no error")

	// v, err = kv.AtomicDelete("test-key-2")
	// assert.Nil(t, v, "should be nil")
	// assert.Equal(t, ErrKeyNotFound, err, "should be equal")
}
