package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

type Response struct {
	resp *http.Response
}

func (r *Response) Code() int {
	return r.resp.StatusCode
}

func (r *Response) Message() string {
	return r.resp.Status
}

func (r *Response) GetHeader(name string) string {
	return r.resp.Header.Get(name)
}

func (r *Response) GetHeaders(name string) []string {
	return r.resp.Header.Values(name)
}

func (r *Response) Body() io.ReadCloser {
	return r.resp.Body
}

type ResponseReader interface {
	ReadResponse(context.Context, *Response) error
}

func (cli *Client) handleRequestWithReader(ctx context.Context, req *http.Request, reader ResponseReader) error {
	// do request
	start := time.Now()
	if !cli.DisableVerboseLog {
		d, _ := httputil.DumpRequestOut(req, true)
		if cli.UseZapLogger {
			logger.Debug(req.Context(), "HTTP client handleRequest",
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
				zap.String("request", string(d)),
			)
		} else {
			logger.Debugf(ctx, "%s %s", req.Method, req.URL.String())
			logger.Debugf(ctx, ">>>>>>>>>>>>>>>: %v", string(d))
		}
	}

	resp, err := cli.HTTPClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if !cli.DisableVerboseLog {
		dr, _ := httputil.DumpResponse(resp, true)
		if cli.UseZapLogger {
			logger.Debug(req.Context(), "HTTP client handleRequest",
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
				zap.String("result", string(dr)),
				zap.Int("statusCode", resp.StatusCode),
				zap.Duration("latency", time.Since(start)),
			)
		} else {
			logger.Debugf(ctx, "\"%s %s\": code=%v, consumed=%v", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))
			logger.Debugf(ctx, ">>>>>>>>>>>>>>>: %v", string(dr))
		}
	}

	return reader.ReadResponse(ctx, &Response{resp})
}

// DoWithReader 发送HTTP请求
func (cli *Client) DoWithReader(ctx context.Context, method string, path string, header http.Header, params url.Values, data interface{}, reader ResponseReader, contentType string) error {
	// prepare request
	var rb io.Reader
	if data != nil {
		if reader, ok := data.(io.Reader); ok {
			rb = reader
		} else {
			content, err := json.Marshal(data)
			if err != nil {
				return err
			}

			rb = bytes.NewReader(content)
		}
	}

	sumParms := make(url.Values)

	reqUrl := cli.Host + path
	if cli.CommonParams != nil {
		for key, values := range cli.CommonParams {
			for _, value := range values {
				sumParms.Add(key, value)
			}
		}
	}

	if params != nil {
		for key, values := range params {
			for _, value := range values {
				sumParms.Add(key, value)
			}
		}
	}

	if len(sumParms) > 0 {
		reqUrl = reqUrl + "?" + sumParms.Encode()
	}

	logger.Debugf(ctx, "url is %+v", reqUrl)

	req, err := http.NewRequestWithContext(ctx, method, reqUrl, rb)
	if err != nil {
		return err
	}

	req.Header = cli.Header.Clone()
	for k, v := range header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Accept", contentType)
	}

	tr := cli.HTTPClient.Transport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return cli.handleRequestWithReader(ctx, req, reader)
}

// DoRequestWithReader 发送HTTP请求
func (cli *Client) DoRequestWithReader(ctx context.Context, method string, path string, header http.Header, params interface{}, data interface{}, reader ResponseReader, contentType string) error {
	if params == nil {
		return cli.DoWithReader(ctx, method, path, header, nil, data, reader, contentType)
	}

	if values, ok := params.(url.Values); ok {
		return cli.DoWithReader(ctx, method, path, header, values, data, reader, contentType)
	}

	values, err := query.Values(params)
	if err != nil {
		return fmt.Errorf("failed to convert query params to url.Values: %v", err)
	}

	return cli.DoWithReader(ctx, method, path, header, values, data, reader, contentType)
}

// GetWithReader 发送GET请求
func (cli *Client) GetWithReader(ctx context.Context, path string, header http.Header, params interface{}, reader ResponseReader, contentType string) error {
	return cli.DoRequestWithReader(ctx, http.MethodGet, path, header, params, nil, reader, contentType)
}

// HeadWithReader 发送Head请求
func (cli *Client) HeadWithReader(ctx context.Context, path string, header http.Header, params interface{}, reader ResponseReader, contentType string) error {
	return cli.DoRequestWithReader(ctx, http.MethodHead, path, header, params, nil, reader, contentType)
}

// PostWithReader 发送POST请求
func (cli *Client) PostWithReader(ctx context.Context, path string, header http.Header, params interface{}, data interface{}, reader ResponseReader, contentType string) error {
	return cli.DoRequestWithReader(ctx, http.MethodPost, path, header, params, data, reader, contentType)
}

// PatchWithReader 发送PATCH请求
func (cli *Client) PatchWithReader(ctx context.Context, path string, header http.Header, params interface{}, data interface{}, reader ResponseReader, contentType string) error {
	return cli.DoRequestWithReader(ctx, http.MethodPatch, path, header, params, data, reader, contentType)
}

// PutWithReader 发送PUT请求
func (cli *Client) PutWithReader(ctx context.Context, path string, header http.Header, params interface{}, data interface{}, reader ResponseReader, contentType string) error {
	return cli.DoRequestWithReader(ctx, http.MethodPut, path, header, params, data, reader, contentType)
}

// DeleteWithReader 发送Delete请求
func (cli *Client) DeleteWithReader(ctx context.Context, path string, header http.Header, params interface{}, reader ResponseReader, contentType string) error {
	return cli.DoRequestWithReader(ctx, http.MethodDelete, path, header, params, nil, reader, contentType)
}
