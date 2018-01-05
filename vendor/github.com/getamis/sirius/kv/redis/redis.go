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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	store "github.com/getamis/sirius/kv"
	"github.com/getamis/sirius/log"
	"gopkg.in/redis.v5"
)

var (
	// ErrMultipleEndpointsUnsupported is thrown when there are
	// multiple endpoints specified for Redis
	ErrMultipleEndpointsUnsupported = errors.New("redis does not support multiple endpoints")

	// ErrAbortTryLock is thrown when a user stops trying to seek the lock
	// by sending a signal to the stop chan, this is used to verify if the
	// operation succeeded
	ErrAbortTryLock = errors.New("redis: lock operation aborted")
)

// New creates a new Redis client given a list
// of endpoints and optional config
func New(endpoints []string, options ...RedisOption) (*Redis, error) {
	if len(endpoints) > 1 {
		return nil, ErrMultipleEndpointsUnsupported
	}

	return new(endpoints, options...), nil
}

func new(endpoints []string, options ...RedisOption) *Redis {
	rOption := &redis.Options{
		Addr:         endpoints[0],
		Password:     "", //TODO: make sure that the password can be supported in *redis.ClusterClient
		DialTimeout:  5 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	for _, option := range options {
		option(rOption)
	}

	return &Redis{
		client: redis.NewClient(rOption),
		script: redis.NewScript(luaScript()),
		codec:  defaultCodec{},
	}
}

type defaultCodec struct{}

func (c defaultCodec) encode(kv *store.KeyValue) (string, error) {
	b, err := json.Marshal(kv)
	return string(b), err
}

func (c defaultCodec) decode(b string, kv *store.KeyValue) error {
	return json.Unmarshal([]byte(b), kv)
}

// Redis implements libkv.Store interface with redis backend
type Redis struct {
	client *redis.Client
	script *redis.Script
	codec  defaultCodec
}

const (
	noExpiration   = time.Duration(0)
	defaultLockTTL = 60 * time.Second
)

// Put a value at the specified key
func (r *Redis) Put(key string, value []byte, options ...store.PutOption) error {
	pOption := &store.PutOptions{}
	for _, option := range options {
		option(pOption)
	}
	expirationAfter := noExpiration
	if pOption.TTL != noExpiration {
		expirationAfter = pOption.TTL
	}

	return r.setTTL(normalize(key), &store.KeyValue{
		Key:      key,
		Value:    value,
		Revision: sequenceNum(),
	}, expirationAfter)
}

func (r *Redis) setTTL(key string, val *store.KeyValue, ttl time.Duration) error {
	valStr, err := r.codec.encode(val)
	if err != nil {
		return err
	}

	return r.client.Set(key, valStr, ttl).Err()
}

// Get a value given its key
func (r *Redis) Get(key string) (*store.KeyValue, error) {
	return r.get(normalize(key))
}

func (r *Redis) get(key string) (*store.KeyValue, error) {
	reply, err := r.client.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, store.ErrKeyNotFound
		}
		return nil, err
	}
	val := store.KeyValue{}
	if err := r.codec.decode(string(reply), &val); err != nil {
		return nil, err
	}
	return &val, nil
}

// Delete the value at the specified key
func (r *Redis) Delete(key string) error {
	return r.client.Del(normalize(key)).Err()
}

// Exists verify if a Key exists in the store
func (r *Redis) Exists(key string) (bool, error) {
	return r.client.Exists(normalize(key)).Result()
}

// Watch for changes on a key
// glitch: we use notified-then-retrieve to retrieve *store.KeyValue.
// so the responses may sometimes inaccurate
func (r *Redis) Watch(key string, stopCh <-chan struct{}) (<-chan *store.KeyValue, error) {
	watchCh := make(chan *store.KeyValue)
	nKey := normalize(key)

	get := getter(func() (interface{}, error) {
		pair, err := r.get(nKey)
		if err != nil {
			return nil, err
		}
		return pair, nil
	})

	push := pusher(func(v interface{}) {
		if val, ok := v.(*store.KeyValue); ok {
			watchCh <- val
		}
	})

	sub, err := newSubscribe(r.client, regexWatch(nKey, false))
	if err != nil {
		return nil, err
	}

	go func(sub *subscribe, stopCh <-chan struct{}, get getter, push pusher) {
		defer sub.Close()

		msgCh := sub.Receive(stopCh)
		if err := watchLoop(msgCh, stopCh, get, push); err != nil {
			log.Info("watchLoop in Watch", "err", err)
		}
	}(sub, stopCh, get, push)

	return watchCh, nil
}

func regexWatch(key string, withChildren bool) string {
	var regex string
	if withChildren {
		regex = fmt.Sprintf("__keyspace*:%s*", key)
		// for all database and keys with $key prefix
	} else {
		regex = fmt.Sprintf("__keyspace*:%s", key)
		// for all database and keys with $key
	}
	return regex
}

// getter defines a func type which retrieves data from remote storage
type getter func() (interface{}, error)

// pusher defines a func type which pushes data blob into watch channel
type pusher func(interface{})

func watchLoop(msgCh chan *redis.Message, stopCh <-chan struct{}, get getter, push pusher) error {

	// deliver the original data before we setup any events
	pair, err := get()
	if err != nil {
		return err
	}
	push(pair)

	for range msgCh {
		// retrieve and send back
		pair, err := get()
		if err != nil {
			return err
		}
		push(pair)
	}

	return nil
}

type subscribe struct {
	pubsub  *redis.PubSub
	closeCh chan struct{}
}

func newSubscribe(client *redis.Client, regex string) (*subscribe, error) {
	ch, err := client.PSubscribe(regex)
	if err != nil {
		return nil, err
	}
	return &subscribe{
		pubsub:  ch,
		closeCh: make(chan struct{}),
	}, nil
}

func (s *subscribe) Close() error {
	close(s.closeCh)
	return s.pubsub.Close()
}

func (s *subscribe) Receive(stopCh <-chan struct{}) chan *redis.Message {
	msgCh := make(chan *redis.Message)
	go s.receiveLoop(msgCh, stopCh)
	return msgCh
}

func (s *subscribe) receiveLoop(msgCh chan *redis.Message, stopCh <-chan struct{}) {
	defer close(msgCh)

	for {
		select {
		case <-s.closeCh:
			return
		case <-stopCh:
			return
		default:
			msg, err := s.pubsub.ReceiveMessage()
			if err != nil {
				return
			}
			if msg != nil {
				msgCh <- msg
			}
		}
	}
}

// WatchTree watches for changes on child nodes under
// a given directory
func (r *Redis) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*store.KeyValue, error) {
	watchCh := make(chan []*store.KeyValue)
	nKey := normalize(directory)

	get := getter(func() (interface{}, error) {
		pair, err := r.list(nKey)
		if err != nil {
			return nil, err
		}
		return pair, nil
	})

	push := pusher(func(v interface{}) {
		if _, ok := v.([]*store.KeyValue); !ok {
			return
		}
		watchCh <- v.([]*store.KeyValue)
	})

	sub, err := newSubscribe(r.client, regexWatch(nKey, true))
	if err != nil {
		return nil, err
	}

	go func(sub *subscribe, stopCh <-chan struct{}, get getter, push pusher) {
		defer sub.Close()

		msgCh := sub.Receive(stopCh)
		if err := watchLoop(msgCh, stopCh, get, push); err != nil {
			log.Info("watchLoop in WatchTree", "err", err)
		}
	}(sub, stopCh, get, push)

	return watchCh, nil
}

// NewLock creates a lock for a given key.
// The returned Locker is not held and must be acquired
// with `.Lock`. The Value is optional.
func (r *Redis) Lock(key string, options ...store.LockOption) (store.Locker, error) {
	var (
		value []byte
		ttl   = defaultLockTTL
	)

	lOption := &store.LockOptions{}
	for _, option := range options {
		option(lOption)
	}

	if lOption.TTL != noExpiration {
		ttl = lOption.TTL
	}
	if len(lOption.Value) != 0 {
		value = lOption.Value
	}

	return &redisLock{
		redis:    r,
		last:     nil,
		key:      key,
		value:    value,
		ttl:      ttl,
		unlockCh: make(chan struct{}),
	}, nil
}

type redisLock struct {
	redis    *Redis
	last     *store.KeyValue
	unlockCh chan struct{}

	key   string
	value []byte
	ttl   time.Duration
}

func (l *redisLock) Lock(stopCh chan struct{}) (<-chan struct{}, error) {
	lockHeld := make(chan struct{})

	success, err := l.tryLock(lockHeld, stopCh)
	if err != nil {
		return nil, err
	}
	if success {
		return lockHeld, nil
	}

	// wait for changes on the key
	watch, err := l.redis.Watch(l.key, stopCh)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case <-stopCh:
			return nil, ErrAbortTryLock
		case <-watch:
			success, err := l.tryLock(lockHeld, stopCh)
			if err != nil {
				return nil, err
			}
			if success {
				return lockHeld, nil
			}
		}
	}
}

// tryLock return true, nil when it acquired and hold the lock
// and return false, nil when it can't lock now,
// and return false, err if any unespected error happened underlying
func (l *redisLock) tryLock(lockHeld, stopChan chan struct{}) (bool, error) {
	success, new, err := l.redis.atomicPut(
		l.key,
		l.value,
		l.last,
		store.PutExpiration(l.ttl),
	)
	if success {
		l.last = new
		// keep holding
		go l.holdLock(lockHeld, stopChan)
		return true, nil
	}
	if err != nil && (err == store.ErrKeyNotFound || err == store.ErrKeyModified || err == store.ErrKeyExists) {
		return false, nil
	}
	return false, err
}

func (l *redisLock) holdLock(lockHeld, stopChan chan struct{}) {
	defer close(lockHeld)

	hold := func() error {
		_, new, err := l.redis.atomicPut(
			l.key,
			l.value,
			l.last,
			store.PutExpiration(l.ttl),
		)
		if err == nil {
			l.last = new
		}
		return err
	}

	heartbeat := time.NewTicker(l.ttl / 3)
	defer heartbeat.Stop()

	for {
		select {
		case <-heartbeat.C:
			if err := hold(); err != nil {
				return
			}
		case <-l.unlockCh:
			return
		case <-stopChan:
			return
		}
	}
}

func (l *redisLock) Unlock() error {
	l.unlockCh <- struct{}{}

	_, err := l.redis.AtomicDelete(l.key, l.last)
	if err != nil {
		return err
	}
	l.last = nil

	return err
}

// List the content of a given prefix
func (r *Redis) List(directory string) ([]*store.KeyValue, error) {
	return r.list(normalize(directory))
}

func (r *Redis) list(directory string) ([]*store.KeyValue, error) {

	var allKeys []string
	regex := scanRegex(directory) // for all keyed with $directory
	allKeys, err := r.keys(regex)
	if err != nil {
		return nil, err
	}
	// TODO: need to handle when #key is too large
	return r.mget(directory, allKeys...)
}

func (r *Redis) keys(regex string) ([]string, error) {
	const (
		startCursor  = 0
		endCursor    = 0
		defaultCount = 10
	)

	var allKeys []string

	keys, nextCursor, err := r.client.Scan(startCursor, regex, defaultCount).Result()
	if err != nil {
		return nil, err
	}
	allKeys = append(allKeys, keys...)
	for nextCursor != endCursor {
		keys, nextCursor, err = r.client.Scan(nextCursor, regex, defaultCount).Result()
		if err != nil {
			return nil, err
		}

		allKeys = append(allKeys, keys...)
	}
	if len(allKeys) == 0 {
		return nil, store.ErrKeyNotFound
	}
	return allKeys, nil
}

// mget values given their keys
func (r *Redis) mget(directory string, keys ...string) ([]*store.KeyValue, error) {
	replies, err := r.client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	pairs := []*store.KeyValue{}
	for _, reply := range replies {
		var sreply string
		if _, ok := reply.(string); ok {
			sreply = reply.(string)
		}
		if sreply == "" {
			// empty reply
			continue
		}

		newkv := &store.KeyValue{}
		if err := r.codec.decode(sreply, newkv); err != nil {
			return nil, err
		}
		if normalize(newkv.Key) != directory {
			pairs = append(pairs, newkv)
		}
	}
	return pairs, nil
}

// DeleteTree deletes a range of keys under a given directory
// glitch: we list all available keys first and then delete them all
// it costs two operations on redis, so is not atomicity.
func (r *Redis) DeleteTree(directory string) error {
	var allKeys []string
	regex := scanRegex(normalize(directory)) // for all keyed with $directory
	allKeys, err := r.keys(regex)
	if err != nil {
		return err
	}
	return r.client.Del(allKeys...).Err()
}

// AtomicPut is an atomic CAS operation on a single value.
// Pass previous = nil to create a new key.
// we introduced script on this page, so atomicity is guaranteed
func (r *Redis) AtomicPut(key string, value []byte, previous *store.KeyValue, options ...store.PutOption) (*store.KeyValue, error) {
	_, result, err := r.atomicPut(key, value, previous, options...)
	return result, err
}

func (r *Redis) atomicPut(key string, value []byte, previous *store.KeyValue, options ...store.PutOption) (bool, *store.KeyValue, error) {
	pOption := &store.PutOptions{}
	for _, option := range options {
		option(pOption)
	}
	expirationAfter := noExpiration
	if pOption.TTL != noExpiration {
		expirationAfter = pOption.TTL
	}

	newKV := &store.KeyValue{
		Key:      key,
		Value:    value,
		Revision: sequenceNum(),
	}
	nKey := normalize(key)

	// if previous == nil, set directly
	if previous == nil {
		if err := r.setNX(nKey, newKV, expirationAfter); err != nil {
			return false, nil, err
		}
		return true, newKV, nil
	}

	if err := r.cas(
		nKey,
		previous,
		newKV,
		formatSec(expirationAfter),
	); err != nil {
		return false, nil, err
	}
	return true, newKV, nil
}

func (r *Redis) setNX(key string, val *store.KeyValue, expirationAfter time.Duration) error {
	valBlob, err := r.codec.encode(val)
	if err != nil {
		return err
	}

	if !r.client.SetNX(key, valBlob, expirationAfter).Val() {
		return store.ErrKeyExists
	}
	return nil
}

func (r *Redis) cas(key string, old, new *store.KeyValue, secInStr string) error {
	newVal, err := r.codec.encode(new)
	if err != nil {
		return err
	}

	oldVal, err := r.codec.encode(old)
	if err != nil {
		return err
	}

	return r.runScript(
		cmdCAS,
		key,
		oldVal,
		newVal,
		secInStr,
	)
}

// AtomicDelete is an atomic delete operation on a single value
// the value will be deleted if previous matched the one stored in db
func (r *Redis) AtomicDelete(key string, previous *store.KeyValue) (bool, error) {
	if err := r.cad(normalize(key), previous); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Redis) cad(key string, old *store.KeyValue) error {
	oldVal, err := r.codec.encode(old)
	if err != nil {
		return err
	}

	return r.runScript(
		cmdCAD,
		key,
		oldVal,
	)
}

// Close the store connection
func (r *Redis) Close() {
	r.client.Close()
}

func scanRegex(directory string) string {
	return fmt.Sprintf("%s*", directory)
}

func (r *Redis) runScript(args ...interface{}) error {
	err := r.script.Run(
		r.client,
		nil,
		args...,
	).Err()
	if err != nil && strings.Contains(err.Error(), "redis: key is not found") {
		return store.ErrKeyNotFound
	}
	if err != nil && strings.Contains(err.Error(), "redis: value has been changed") {
		return store.ErrKeyModified
	}
	return err
}

func normalize(key string) string {
	return store.Normalize(key)
}

func formatSec(dur time.Duration) string {
	return fmt.Sprintf("%d", int(dur/time.Second))
}

func sequenceNum() uint64 {
	// TODO: use uuid if we concerns collision probability of this number
	return uint64(time.Now().Nanosecond())
}
