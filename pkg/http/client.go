package client

import (
	"bytes"
	"crypto/tls"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
	klog "k8s.io/klog/v2"
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

// Config 服务Client的配置
type Config struct {
	Host              string `json:"host" mapstructure:"host"`
	UserName          string
	UserPassword      string
	TokenTimeOut      time.Duration
	HTTP              HTTPConfig `mapstructure:"http"`
	DisableVerboseLog bool       `json:"disableVerboseLog" mapstructure:"disableVerboseLog"`
	UseZapLogger      bool       `json:"useZapLogger" mapstructure:"useZapLogger"`
}

type httpError struct {
	HTTPCode int    `json:"httpCode"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Err      string `json:"error"`
}

// HTTPError 请求错误响应
type HTTPError struct {
	httpError
	Status string `json:"status"`
}

func (e *HTTPError) Error() string {
	detail := e.Err
	if len(detail) == 0 {
		detail = e.Status
	}
	return fmt.Sprintf("HTTPCode=%d, %s(%s): %s", e.HTTPCode, e.Message, e.Code, detail)
}

// Client HTTP客户端
type Client struct {
	ProtocolHostPort  string
	UserName          string
	UserPassword      string
	ContentType       string
	LoginPath         string
	TopoPath          string
	DisableVerboseLog bool
	UseZapLogger      bool
	Header            http.Header
	Token             string
	TokenTimeOutAt    time.Time
	TokenTimeOut      time.Duration
	HTTPClient        *http.Client
}

var CCAEClient *Client

// NewClient 返回一个新的客户端
func NewClient(conf *config.ServerConfig) *Client {
	protocolHostPort := fmt.Sprintf("%s://%s:%s", conf.CCAE.Server.Protocol, conf.CCAE.Server.Host, conf.CCAE.Server.Port)
	CCAEClient = &Client{
		ProtocolHostPort:  protocolHostPort,
		LoginPath:         conf.CCAE.Api.TokenUrl,
		TopoPath:          conf.CCAE.Api.QureyTopoUrl,
		UserName:          conf.CCAE.Server.UserName,
		UserPassword:      conf.CCAE.Server.UserPassword,
		TokenTimeOut:      time.Duration(conf.CCAE.Server.TokenTimeOut),
		ContentType:       "application/json",
		Header:            make(http.Header),
		HTTPClient:        NewHTTPClient(&HTTPConfig{}),
		DisableVerboseLog: true,
		UseZapLogger:      true,
	}
	return CCAEClient
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

	// caCert, err := ioutil.ReadFile("path/to/ca.pem")
	// if err != nil {
	// 	panic(err)
	// }
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// tlsConfig := &tls.Config{
	// 	RootCAs: caCertPool,
	// }

	return &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dial.Dial,
			DialContext:         dial.DialContext,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			//TLSClientConfig:     tlsConfig,
		},
		Timeout: time.Duration(timeout) * time.Millisecond,
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

	url := cli.ProtocolHostPort + path
	req, err := http.NewRequest(method, url, rb)
	if err != nil {
		return nil, err
	}
	req.Header = cli.Header.Clone()
	req.Header.Add("X-Auth-Token", cli.Token)
	req.Header.Add("Content-Type", cli.ContentType)
	req.URL.RawQuery = params.Encode()

	// do request
	return cli.handleRequest(req)
}

func (cli *Client) handleRequest(req *http.Request) ([]byte, error) {
	// do request
	start := time.Now()
	resp, err := cli.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !cli.DisableVerboseLog {
		if cli.UseZapLogger {
			klog.Info(req.Context(), "HTTP client handleRequest",
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
				zap.Int("statusCode", resp.StatusCode),
				zap.Duration("latency", time.Since(start)))
		} else {
			klog.Infof("\"%s %s\": code=%v, consumed=%v", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))
		}
	}

	// handle response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	klog.Infof("resp.Body:%+v", string(content))
	klog.Infof("resp.StatusCode:%+v", resp.StatusCode)

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		// var errorInfo HTTPError
		// if err := json.Unmarshal(content, &errorInfo); err == nil {
		// 	errorInfo.HTTPCode = resp.StatusCode
		// 	if errorInfo.Code == "" && errorInfo.Status != "" {
		// 		errorInfo.Code = errorInfo.Status
		// 	}
		// 	return nil, &errorInfo
		// } else if err := json.Unmarshal(content, &errorInfo.httpError); err == nil {
		// 	errorInfo.HTTPCode = resp.StatusCode
		// 	return nil, &errorInfo
		// }
		var errorInfo common.ErrorInfoResp
		err := json.Unmarshal(content, &errorInfo)
		if err != nil {
			return nil, fmt.Errorf("%s", string(content))
		}
		if errorInfo.DetailArgs == nil {
			return nil, fmt.Errorf("%s", errorInfo.DetailArgs)
		}
		return nil, fmt.Errorf("%s", errorInfo.DetailArgs)
	}
	return content, nil
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

func Format(val interface{}) (string, error) {
	var err error
	switch v := val.(type) {
	case string:
		return v, err
	case int:
		return strconv.Itoa(v), err
	case int32:
		return strconv.FormatInt(int64(v), 10), err
	case int64:
		return strconv.FormatInt(v, 10), err
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), err
	case bool:
		return strconv.FormatBool(v), err
	default:
		return fmt.Sprintf("%v", v), err
	}
}

func Convert(params interface{}) (url.Values, error) {
	if values, ok := params.(url.Values); ok {
		return values, nil
	}
	content, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var maps map[string]interface{}
	if err = json.Unmarshal(content, &maps); err != nil {
		return nil, err
	}
	values := url.Values{}
	for key, val := range maps {
		var value string
		switch t := val.(type) {
		case int, int32, int64, float32, float64, string, bool:
			str, err := Format(t)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %T(%v=%v) to string: %v", val, key, val, err)
			}
			value = str
		default:
			return nil, fmt.Errorf("unexpected type %T while converting: %v=%v", t, key, val)
		}
		values.Add(key, value)
	}
	return values, nil
}
