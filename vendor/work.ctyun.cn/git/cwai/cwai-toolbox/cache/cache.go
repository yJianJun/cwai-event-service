package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Policy 缓存策略
type Policy string

// WrappedFunc 被缓存结果的函数定义
type WrappedFunc func() ([]byte, error)

const (
	// PolicyNo 不使用缓存
	PolicyNo Policy = ""
	// PolicyWriteOnly 仅记录缓存
	PolicyWriteOnly Policy = "writeOnly"
	// PolicyReadWrite 写缓存并从缓存中读取
	PolicyReadWrite Policy = "readWrite"
)

// Cache 缓存
type Cache struct {
	Directory string
	Policy    Policy
}

// Wrap 缓存函数结果
func (cache *Cache) Wrap(filename string, f WrappedFunc) (content []byte, err error) {
	if cache.Policy == PolicyNo {
		return f()
	}

	if cache.Policy == PolicyWriteOnly {
		return cache.WrapWriteOnly(filename, f)
	}

	if cache.Policy == PolicyReadWrite {
		return cache.WrapReadWrite(filename, f)
	}
	return nil, fmt.Errorf("unknown cache policy %s", cache.Policy)
}

// WrapWriteOnly 仅记录缓存
func (cache *Cache) WrapWriteOnly(filename string, f WrappedFunc) (content []byte, err error) {
	content, err = f()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(cache.Directory, filename)
	cache.EnsureDir(filepath.Dir(path))
	if err = ioutil.WriteFile(path, content, 0644); err != nil {
		return nil, err
	}
	return content, nil
}

// WrapReadWrite 写缓存，并在缓存存在的情况下从缓存读
func (cache *Cache) WrapReadWrite(filename string, f WrappedFunc) (content []byte, err error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return cache.WrapWriteOnly(filename, f)
	}

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// EnsureDir 确保目录存在
func (cache *Cache) EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
