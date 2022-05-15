package main

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, v interface{})
	Get(key string) (interface{}, bool)
}

type TimeoutCache struct {
	mu        sync.Mutex
	timeout   time.Duration
	caches    map[string]interface{}
	beginTime map[string]int64
}

func (tc *TimeoutCache) setTimeout(timeout time.Duration) {
	tc.timeout = timeout
}

func (tc *TimeoutCache) Set(key string, v interface{}) {
	tc.mu.Lock()
	tc.caches[key] = v
	tc.beginTime[key] = time.Now().Unix()
	tc.mu.Unlock()
}

func (tc *TimeoutCache) remove(key string) {
	tc.mu.Lock()
	delete(tc.caches, key)
	delete(tc.beginTime, key)
	tc.mu.Unlock()
}

func (tc *TimeoutCache) Get(key string) (interface{}, bool) {
	tc.mu.Lock()
	if tc.caches[key] == nil {
		return nil, false
	}

	if tc.timeout < time.Duration(time.Now().Unix()-tc.beginTime[key]) {
		tc.remove(key)
		return nil, false
	}

	tc.beginTime[key] = time.Now().Unix()
	tc.mu.Unlock()

	return tc.caches[key], true
}

func NewTimeoutCache(timeout time.Duration) Cache {
	var tc = &TimeoutCache{sync.Mutex{}, timeout, make(map[string]interface{}), make(map[string]int64)}
	tc.setTimeout(timeout)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			tc.mu.Lock()
			for k, _ := range tc.caches {
				if timeout < time.Duration(time.Now().Unix()-tc.beginTime[k]) {
					tc.remove(k)
				}
			}
			tc.mu.Unlock()
		}
	}()
	return tc
}
