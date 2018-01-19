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
	"time"

	"gopkg.in/redis.v5"
)

type RedisOption func(*redis.Options)

func DialTimeout(t time.Duration) RedisOption {
	return func(r *redis.Options) {
		r.DialTimeout = t
	}
}

func ReadTimeout(t time.Duration) RedisOption {
	return func(r *redis.Options) {
		r.ReadTimeout = t
	}
}

func WriteTimeout(t time.Duration) RedisOption {
	return func(r *redis.Options) {
		r.WriteTimeout = t
	}
}

func Password(pwd string) RedisOption {
	return func(r *redis.Options) {
		r.Password = pwd
	}
}
