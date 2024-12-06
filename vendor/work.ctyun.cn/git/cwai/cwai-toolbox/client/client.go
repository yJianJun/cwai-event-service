package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Config 服务Client的配置
type Config struct {
	Host              string     `json:"host" mapstructure:"host"`
	Token             string     `json:"token" mapstructure:"token"`
	HTTP              HTTPConfig `mapstructure:"http"`
	DisableVerboseLog bool       `json:"disableVerboseLog" mapstructure:"disableVerboseLog"`
	UseZapLogger      bool       `json:"useZapLogger" mapstructure:"useZapLogger"`
	CommonParams      url.Values `json:"commonParams" mapstructure:"commonParams"`
}

// ConfigOption 配置选项
type ConfigOption func(c *Config)

// WithTimeout 设置请求超时时间 ms
func WithTimeout(d int) ConfigOption {
	return func(c *Config) {
		c.HTTP.Timeout = d
	}
}

// WithDialTimeout 设置连接超时时间 ms
func WithDialTimeout(d int) ConfigOption {
	return func(c *Config) {
		c.HTTP.DialTimeout = d
	}
}

// WithCustomConfig 自定义请求设置
func WithCustomConfig(host, token string, opts ...ConfigOption) *Config {
	c := &Config{
		Host:  host,
		Token: token,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// Client HTTP客户端
type Client struct {
	Host              string
	Token             string
	ContentType       string
	DisableVerboseLog bool
	UseZapLogger      bool
	Header            http.Header
	HTTPClient        *http.Client
	CommonParams      url.Values
}

// NewClient 返回一个新的客户端
func NewClient(conf *Config) *Client {
	return &Client{
		Host:              conf.Host,
		Token:             conf.Token,
		ContentType:       "application/json",
		Header:            make(http.Header),
		HTTPClient:        NewHTTPClient(&conf.HTTP),
		DisableVerboseLog: conf.DisableVerboseLog,
		UseZapLogger:      conf.UseZapLogger,
		CommonParams:      conf.CommonParams,
	}
}

// Clone 生成一个新的Client，但是会共享同一个HTTP Client
func (cli *Client) Clone() *Client {
	newCli := *cli
	newCli.Header = cli.Header.Clone()
	return &newCli
}

// WithHeader 附加额外的Header到后续的HTTP请求中
func (cli *Client) WithHeader(key, value string) *Client {
	newCli := cli.Clone()
	newCli.Header.Add(key, value)
	return newCli
}

// WithoutVerboseLog 设置不打印详细日志
func (cli *Client) WithoutVerboseLog() *Client {
	newCli := cli.Clone()
	newCli.DisableVerboseLog = true
	return newCli
}

// Do 发送HTTP请求
func (cli *Client) Do(method string, path string, pid int, params url.Values, data interface{}) ([]byte, error) {
	// prepare request
	var rb io.Reader
	if data != nil {
		if reader, ok := data.(io.Reader); ok {
			rb = reader
		} else {
			content, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}

			rb = bytes.NewReader(content)
		}
	}

	url := cli.Host + path
	req, err := http.NewRequest(method, url, rb)
	if err != nil {
		return nil, err
	}
	req.Header = cli.Header.Clone()
	req.Header.Add("cwaiToken", cli.Token) // todo toolbox是个业务无关的工程，不能依赖api工程，所以这里暂时用字符串代替，更好的办法是在pai工程里封装这个方法
	req.Header.Add("Content-Type", cli.ContentType)
	req.URL.RawQuery = params.Encode()

	// do request
	return cli.handleRequest(req)
}

// DoRequest 发送HTTP请求
func (cli *Client) DoRequest(method string, path string, pid int, params interface{}, data interface{}) ([]byte, error) {
	if params == nil {
		return cli.Do(method, path, pid, nil, data)
	}
	if values, ok := params.(url.Values); ok {
		return cli.Do(method, path, pid, values, data)
	}
	values, err := Convert(params)
	if err != nil {
		return nil, fmt.Errorf("failed to convert query params to url.Values: %v", err)
	}
	return cli.Do(method, path, pid, values, data)
}

/*
examples #1:
	cli.Get("/url", 0, nil)

examples #2:
	params := map[string]string{"key": "value"}
	cli.Get("/url", 0, params)

examples #3:
	type Query struct {
		Key string `json:"key,omitempty"`
	}
	params := Query{Key: "value"}
	cli.Get("/url", 0, params)
*/
// Get 发送GET请求
func (cli *Client) Get(path string, projectID int, params interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodGet, path, projectID, params, nil)
}

// Post 发送POST请求
func (cli *Client) Post(path string, projectID int, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodPost, path, projectID, params, data)
}

// Patch 发送PATCH请求
func (cli *Client) Patch(path string, projectID int, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodPatch, path, projectID, params, data)
}

// Put 发送PUT请求
func (cli *Client) Put(path string, projectID int, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodPut, path, projectID, params, data)
}

// Delete 发送Delete请求
func (cli *Client) Delete(path string, projectID int, params interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodDelete, path, projectID, params, nil)
}
