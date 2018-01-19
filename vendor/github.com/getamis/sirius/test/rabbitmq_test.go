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

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestRabbitMQContainer(t *testing.T) {
	container, _ := NewRabbitMQContainer()
	assert.NotNil(t, container)
	assert.NoError(t, container.Start())

	conn, err := amqp.Dial(container.URL)
	assert.NoError(t, err, "should be no error")
	conn.Close()

	// stop rabbitmq
	assert.NoError(t, container.Suspend())
	time.Sleep(100 * time.Millisecond)
	_, err = amqp.Dial(container.URL)
	assert.Error(t, err, "should got error")

	// restart rabbitmq
	assert.NoError(t, container.Start())
	time.Sleep(time.Second)
	conn, err = amqp.Dial(container.URL)
	assert.NoError(t, err, "should be no error")
	conn.Close()

	// close rabbitmq
	container.Stop()
	time.Sleep(100 * time.Millisecond)
	_, err = amqp.Dial(container.URL)
	assert.Error(t, err, "should got error")
}
