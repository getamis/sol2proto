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

import (
	"crypto/tls"
)

type Options struct {
	Secure    bool
	TLSConfig *tls.Config
}

// Secure communication with the broker
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig specifies the TLS config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

type Option func(*Options)

// ----------------------------------------------------------------------------

type PublishOptions struct {
}

type PublishOption func(*PublishOptions)

// ----------------------------------------------------------------------------

type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string
}

type SubscribeOption func(*SubscribeOptions)

// AutoAck will enable/disable auto acking of messages
// after they have been handled.
func AutoAck(enabled bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = enabled
	}
}

// Queue sets the name of the queue to share messages on
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}
