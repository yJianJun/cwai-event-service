package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

func (cli *Client) handleRequest(req *http.Request) ([]byte, error) {
	// do request
	start := time.Now()
	traceID := uuid.New().String()
	if os.Getenv("API_SDK_DEBUG") == "true" {
		requestDump, _ := httputil.DumpRequest(req, true)
		fmt.Printf("cwaiTraceID: %v, request: %v\n", traceID, string(requestDump))
	}
	resp, err := cli.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("API_SDK_DEBUG") == "true" {
		responseDump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("cwaiTraceID: %v, request: %v\n", traceID, string(responseDump))
	}
	defer resp.Body.Close()
	if !cli.DisableVerboseLog {
		if cli.UseZapLogger {
			logger.Info(req.Context(), "HTTP client handleRequest",
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
				zap.Int("statusCode", resp.StatusCode),
				zap.Duration("latency", time.Since(start)),
			)
		} else {
			glog.Infof("\"%s %s\": code=%v, consumed=%v", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))
		}
	}

	// handle response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		var errorInfo HTTPError
		if err := json.Unmarshal(content, &errorInfo); err == nil {
			errorInfo.HTTPCode = resp.StatusCode
			if errorInfo.Code == "" && errorInfo.Status != "" {
				errorInfo.Code = errorInfo.Status
			}
			return nil, &errorInfo
		} else if err := json.Unmarshal(content, &errorInfo.httpError); err == nil {
			errorInfo.HTTPCode = resp.StatusCode
			return nil, &errorInfo
		}
		return nil, fmt.Errorf("[%d] %s", resp.StatusCode, content[:50])
	}
	return content, nil
}

func (cli *Client) handleApiRequest(req *http.Request) ([]byte, string, error) {
	msg := ""
	// do request
	start := time.Now()
	traceID := uuid.New().String()
	requestDump, _ := httputil.DumpRequest(req, true)
	msg += fmt.Sprintln(string(requestDump))
	if os.Getenv("API_SDK_DEBUG") == "true" {
		fmt.Printf("cwaiTraceID: %v, request: %v\n", traceID, string(requestDump))
	}
	resp, err := cli.HTTPClient.Do(req)
	if err != nil {
		msg += fmt.Sprintf("do http request error %s\n", err.Error())
		return nil, msg, err
	}
	responseDump, _ := httputil.DumpResponse(resp, true)
	msg += fmt.Sprintln(string(responseDump))
	if os.Getenv("API_SDK_DEBUG") == "true" {
		fmt.Printf("cwaiTraceID: %v, request: %v\n", traceID, string(responseDump))
	}
	defer resp.Body.Close()
	if !cli.DisableVerboseLog {
		if cli.UseZapLogger {
			logger.Info(req.Context(), "HTTP client handleRequest",
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
				zap.Int("statusCode", resp.StatusCode),
				zap.Duration("latency", time.Since(start)),
			)
		} else {
			glog.Infof("\"%s %s\": code=%v, consumed=%v", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))
		}
	}

	// handle response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg += fmt.Sprintf("read http resp body error %s\n", err.Error())
		return nil, msg, err
	}
	return content, msg, nil
}

// DoContext 发送HTTP请求
func (cli *Client) DoContext(ctx context.Context, method string, path string, header http.Header, params url.Values, data interface{}) ([]byte, error) {
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
	req, err := http.NewRequestWithContext(ctx, method, url, rb)
	if err != nil {
		return nil, err
	}

	req.Header = cli.Header.Clone()
	for k, v := range header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	req.URL.RawQuery = params.Encode()

	return cli.handleRequest(req)
}

// DoRequestContext 发送HTTP请求
func (cli *Client) DoRequestContext(ctx context.Context, method string, path string, header http.Header, params interface{}, data interface{}) ([]byte, error) {
	if params == nil {
		return cli.DoContext(ctx, method, path, header, nil, data)
	}
	if values, ok := params.(url.Values); ok {
		return cli.DoContext(ctx, method, path, header, values, data)
	}
	values, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("failed to convert query params to url.Values: %v", err)
	}

	return cli.DoContext(ctx, method, path, header, values, data)
}

// GetContext 发送GET请求
func (cli *Client) GetContext(ctx context.Context, path string, header http.Header, params interface{}) ([]byte, error) {
	return cli.DoRequestContext(ctx, http.MethodGet, path, header, params, nil)
}

// PostContext 发送POST请求
func (cli *Client) PostContext(ctx context.Context, path string, header http.Header, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequestContext(ctx, http.MethodPost, path, header, params, data)
}

// PatchContext 发送PATCH请求
func (cli *Client) PatchContext(ctx context.Context, path string, header http.Header, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequestContext(ctx, http.MethodPatch, path, header, params, data)
}

// PutContext 发送PUT请求
func (cli *Client) PutContext(ctx context.Context, path string, header http.Header, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequestContext(ctx, http.MethodPut, path, header, params, data)
}

// DeleteContext 发送Delete请求
func (cli *Client) DeleteContext(ctx context.Context, path string, header http.Header, params interface{}) ([]byte, error) {
	return cli.DoRequestContext(ctx, http.MethodDelete, path, header, params, nil)
}
