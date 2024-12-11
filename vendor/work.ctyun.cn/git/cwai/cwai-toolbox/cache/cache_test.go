package cache

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := Cache{
		Directory: ".",
		Policy:    PolicyNo,
	}
	content, err := cache.Wrap("cache.test", func() ([]byte, error) {
		return []byte("abcd"), nil
	})
	assert.Nil(t, err)
	assert.Equal(t, []byte("abcd"), content)

	cache.Policy = PolicyWriteOnly
	content, err = cache.Wrap("cache.test", func() ([]byte, error) {
		return []byte("abcd"), nil
	})
	assert.Nil(t, err)
	assert.Equal(t, []byte("abcd"), content)
	content, err = ioutil.ReadFile("cache.test")
	assert.Nil(t, err)
	assert.Equal(t, []byte("abcd"), content)
}
