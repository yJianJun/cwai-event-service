package client

import (
	"net"
	"net/http"
	"time"
)

const (
	// DefaultDialTimeout 默认的HTTP连接超时
	DefaultDialTimeout = 3000 // 3s
	// DefaultTimeout 默认的HTTP请求超时
	DefaultTimeout = 60000 // 60s
	// DefaultKeepAlive 默认的KeepAlive时间
	DefaultKeepAlive = 60000 // 60s
	// DefaultMaxIdleConns 默认的最大空闲连接数
	DefaultMaxIdleConns = 1000
	// DefaultMaxIdleConnsPerHost 默认的与每台机器的最大空闲连接数
	DefaultMaxIdleConnsPerHost = 100
)

// HTTPConfig HTTP配置
type HTTPConfig struct {
	// DialTimeout 连接超时时间 ms，默认3s
	DialTimeout int `json:"dialTimeout" mapstructure:"dialTimeout"`
	// Timeout 请求超时时间 ms，默认3s
	Timeout             int `json:"timeout" mapstructure:"timeout"`
	KeepAlive           int `json:"keepAlive" mapstructure:"keepAlive"`
	MaxIdleConns        int `json:"maxIdleConns" mapstructure:"maxIdleConns"`
	MaxIdleConnsPerHost int `json:"maxIdleConnsPerHost" mapstructure:"maxIdleConnsPerHost"`
}

// NewHTTPClient 返回一个原生的http.Client
func NewHTTPClient(conf *HTTPConfig) *http.Client {
	dialTimeout := conf.DialTimeout
	if dialTimeout <= 0 {
		dialTimeout = DefaultDialTimeout
	}

	timeout := conf.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	keepAlive := conf.KeepAlive
	if keepAlive <= 0 {
		keepAlive = DefaultKeepAlive
	}

	maxIdleConns := conf.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = DefaultMaxIdleConns
	}
	maxIdleConnsPerHost := conf.MaxIdleConnsPerHost
	if maxIdleConnsPerHost <= 0 {
		maxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}

	dial := &net.Dialer{
		Timeout:   time.Duration(dialTimeout) * time.Millisecond,
		KeepAlive: time.Duration(keepAlive) * time.Millisecond,
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dial.Dial,
			DialContext:         dial.DialContext,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
		},
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
}
