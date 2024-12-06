package utils

import (
	"sync"
	"time"

)

// SafeMap 是一个并发安全的 map，支持查询不到时重新获取数据、超时控制（0永久）以及获取函数定制
type SafeMap struct {
	sync.Mutex
	data       map[string]*Item
	fetchFunc  func(key string) (interface{}, error)
	defaultTTL time.Duration
}

// Item 表示缓存中的一个项
type Item struct {
	value      interface{}
	expiration time.Time
}

// NewSafeMap 创建一个新的 SafeMap
func NewSafeMap(fetchFunc func(key string) (interface{}, error), defaultTTL time.Duration) *SafeMap {
	return &SafeMap{
		data:       make(map[string]*Item),
		fetchFunc:  fetchFunc,
		defaultTTL: defaultTTL,
	}
}

// Get 从缓存中获取数据，如果数据不存在或已过期，则调用 fetchFunc 获取数据
func (sm *SafeMap) Get(key string) (interface{}, error) {
	sm.Lock()
	defer sm.Unlock()

	if item, ok := sm.data[key]; ok {
		if item.expiration.IsZero() || time.Now().Before(item.expiration) {
			return item.value, nil
		}
		// 数据已过期，删除旧数据
		delete(sm.data, key)
	}

	// 数据不存在或已过期，调用 fetchFunc 获取数据
	value, err := sm.fetchFunc(key)
	if err != nil {
		return nil, err
	}

	// 将新数据存入缓存
	sm.data[key] = &Item{
		value:      value,
		expiration: time.Now().Add(sm.defaultTTL),
	}

	return value, nil
}

// Set 手动设置缓存数据，支持永久
func (sm *SafeMap) Set(key string, value interface{}, ttl time.Duration) {
	sm.Lock()
	defer sm.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	sm.data[key] = &Item{
		value:      value,
		expiration: expiration,
	}
}

// SetPermanent 手动设置永久缓存数据
func (sm *SafeMap) SetPermanent(key string, value interface{}) {
	sm.Set(key, value, 0)
}
