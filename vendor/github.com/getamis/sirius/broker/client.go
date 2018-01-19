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

package broker

// Client is a client interface for various message brokers,
// e.g., RabbitMQ, Kafka, NATS, etc.
type Client interface {
	Connect() error
	Disconnect() error
	Publish(string, *Message, ...PublishOption) error
	Subscribe(string, ...SubscribeOption) (Subscription, error)
}

type Message struct {
	Header map[string]string
	Body   []byte
}

// Publication is given to a subscriber for processing
type Publication interface {
	Topic() string
	Message() *Message
	Ack() error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscription interface {
	Topic() string
	Chan() <-chan Publication
	Unsubscribe() error
}

// ----------------------------------------------------------------------------

func NewClient(opts ...Option) Client {
	return nil
}
