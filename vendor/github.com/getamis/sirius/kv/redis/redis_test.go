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

package redis

import (
	"testing"

	store "github.com/getamis/sirius/kv"
	testutils "github.com/getamis/sirius/kv/test"
	"github.com/getamis/sirius/test"
	"github.com/stretchr/testify/assert"
)

var (
	client = "localhost:6379"
)

func makeRedisClient() store.Store {
	kv := new([]string{client})

	// NOTE: please turn on redis's notification
	// before you using watch/watchtree/lock related features
	kv.client.ConfigSet("notify-keyspace-events", "KA")

	return kv
}

func TestRedisStore(t *testing.T) {
	container, err := test.NewRedisContainer()
	if err != nil {
		t.Fatal("failed to new REDIS container")
	}
	assert.NoError(t, container.Start())
	defer container.Stop()
	kv := makeRedisClient()
	lockTTL := makeRedisClient()
	kvTTL := makeRedisClient()

	testutils.RunTestCommon(t, kv)
	testutils.RunTestAtomic(t, kv)
	testutils.RunTestWatch(t, kv)
	testutils.RunTestLock(t, kv)
	testutils.RunTestLockTTL(t, kv, lockTTL)
	testutils.RunTestTTL(t, kv, kvTTL)
	testutils.RunCleanup(t, kv)
}
