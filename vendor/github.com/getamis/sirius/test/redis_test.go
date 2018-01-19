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

package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v5"
)

func TestRedisContainer(t *testing.T) {
	container, _ := NewRedisContainer()
	assert.NotNil(t, container)
	assert.NoError(t, container.Start())

	conn := redis.NewClient(&redis.Options{
		Addr: container.URL,
	})
	assert.NotNil(t, conn)
	assert.NoError(t, conn.Ping().Err(), "should be no error")

	// stop rabbitmq
	assert.NoError(t, container.Suspend())
	time.Sleep(100 * time.Millisecond)
	conn = redis.NewClient(&redis.Options{
		Addr: container.URL,
	})
	assert.Error(t, conn.Ping().Err(), "should got error")

	// restart rabbitmq
	assert.NoError(t, container.Start())
	time.Sleep(time.Second)
	conn = redis.NewClient(&redis.Options{
		Addr: container.URL,
	})
	assert.NotNil(t, conn)
	assert.NoError(t, conn.Ping().Err(), "should be no error")

	// close rabbitmq
	assert.NoError(t, container.Stop())
	time.Sleep(100 * time.Millisecond)
	conn = redis.NewClient(&redis.Options{
		Addr: container.URL,
	})
	assert.Error(t, conn.Ping().Err(), "should got error")
}
