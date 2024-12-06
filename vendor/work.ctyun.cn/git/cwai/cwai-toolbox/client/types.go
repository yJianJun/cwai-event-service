package client

import (
	"fmt"
)

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

func NewHTTPError(httpCode int, code, message, err string) HTTPError {
	return HTTPError{
		httpError: httpError{
			HTTPCode: httpCode,
			Code:     code,
			Message:  message,
			Err:      err,
		},
	}
}
