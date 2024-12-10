package utils

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 模拟获取数据的函数
func mockFetchFunc(key string) (interface{}, error) {
	return "value for " + key, nil
}

func TestSafeMap_Get(t *testing.T) {
	// 创建一个 SafeMap，设置默认 TTL 为 1 秒
	safeMap := NewSafeMap(mockFetchFunc, 1*time.Second)

	// 测试获取不存在的数据
	value, err := safeMap.Get("key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "value for key1" {
		t.Errorf("Expected 'value for key1', got %v", value)
	}

	// 测试获取已缓存的数据
	value, err = safeMap.Get("key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "value for key1" {
		t.Errorf("Expected 'value for key1', got %v", value)
	}

	// 测试数据过期
	time.Sleep(2 * time.Second)
	value, err = safeMap.Get("key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "value for key1" {
		t.Errorf("Expected 'value for key1', got %v", value)
	}
}

func TestSafeMap_Set(t *testing.T) {
	// 创建一个 SafeMap，设置默认 TTL 为 1 秒
	safeMap := NewSafeMap(mockFetchFunc, 1*time.Second)

	// 测试手动设置缓存数据
	safeMap.Set("key2", "manual value", 2*time.Second)
	value, err := safeMap.Get("key2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "manual value" {
		t.Errorf("Expected 'manual value', got %v", value)
	}

	// 测试手动设置的缓存数据过期
	time.Sleep(3 * time.Second)
	value, err = safeMap.Get("key2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "value for key2" {
		t.Errorf("Expected 'value for key2', got %v", value)
	}
}

func TestSafeMap_SetPermanent(t *testing.T) {
	// 创建一个 SafeMap，设置默认 TTL 为 1 秒
	safeMap := NewSafeMap(mockFetchFunc, 1*time.Second)

	// 测试手动设置永久缓存数据
	safeMap.SetPermanent("key3", "permanent value")
	value, err := safeMap.Get("key3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "permanent value" {
		t.Errorf("Expected 'permanent value', got %v", value)
	}

	// 测试永久缓存数据不会过期
	time.Sleep(2 * time.Second)
	value, err = safeMap.Get("key3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "permanent value" {
		t.Errorf("Expected 'permanent value', got %v", value)
	}
}

func TestSafeMap_FetchFuncError(t *testing.T) {
	// 创建一个 SafeMap，设置默认 TTL 为 1 秒
	safeMap := NewSafeMap(func(key string) (interface{}, error) {
		return nil, errors.New("fetch error")
	}, 1*time.Second)

	// 测试获取数据时 fetchFunc 返回错误
	value, err := safeMap.Get("key4")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if value != nil {
		t.Errorf("Expected nil, got %v", value)
	}
}

func TestSafeMap_ConcurrentAccess(t *testing.T) {
	// 创建一个 SafeMap，设置默认 TTL 为 1 秒
	safeMap := NewSafeMap(mockFetchFunc, 1*time.Second)

	// 并发访问测试
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10)
			value, err := safeMap.Get(key)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			expectedValue := "value for " + key
			if value != expectedValue {
				t.Errorf("Expected '%s', got %v", expectedValue, value)
			}
		}(i)
	}
	wg.Wait()
}
