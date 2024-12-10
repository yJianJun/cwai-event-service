package server

import (
	"fmt"
)

type ServiceError struct {
	Code    *ErrorCode
	Message string
}

func NewServiceError(code *ErrorCode, msgAndArgs ...interface{}) ServiceError {
	message := BuildMessage(msgAndArgs...)
	return ServiceError{
		Code:    code,
		Message: message,
	}
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code.String(), e.Message)
}
