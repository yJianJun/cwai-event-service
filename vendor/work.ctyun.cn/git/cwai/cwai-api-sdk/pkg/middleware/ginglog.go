package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	apicommon "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

func setupLogging(duration time.Duration) {
	go func() {
		for range time.Tick(duration) {
			logger.Flush()
		}
	}()
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// ErrorLogger returns an ErrorLoggerT with parameter gin.ErrorTypeAny
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns an ErrorLoggerT middleware with the given
// type gin.ErrorType.
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(typ).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

// Example:
//
//	router := gin.New()
//	router.Use(ginglog.Logger(3 * time.Second))
func Logger(duration time.Duration) gin.HandlerFunc {
	setupLogging(duration)
	return func(c *gin.Context) {
		logger.SetLogFieldsForGinContext(c)
		t := time.Now()
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		path := c.Request.URL.Path

		var message string
		if errorCode, ok := c.Get(apicommon.LogErrorCode); ok {
			message = fmt.Sprintf("GIN|%s%3d%s|%v|%s|%s %#v %s%s %s%s",
				statusColor, statusCode, reset,
				latency,
				clientIP,
				method,
				path,
				red, errorCode,
				c.Errors.String(), reset,
			)
			logger.Warn(c, message)
		} else {
			message = fmt.Sprintf("GIN|%s%3d%s|%v|%s|%s %#v",
				statusColor, statusCode, reset,
				latency,
				clientIP,
				method,
				path,
			)
			logger.Info(c, message)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
