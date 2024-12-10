package middleware

import (
	"time"

	"go.uber.org/zap"
	zlogger "work.ctyun.cn/git/cwai/cwai-toolbox/logger"

	"github.com/gin-gonic/gin"
)

func setupZapLogging(duration time.Duration) {
	go func() {
		for range time.Tick(duration) {
			zlogger.Flush()
		}
	}()
}

// Logger middleware
// Example:
//
//	router := gin.New()
//	router.Use(middleware.ZapLogger(3 * time.Second))
func ZapLogger(duration time.Duration) gin.HandlerFunc {
	setupZapLogging(duration)
	return ZapLoggerHandler()
}

func ZapLoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		errors := c.Errors

		fields := []zap.Field{zap.Int("statusCode", statusCode),
			zap.Duration("latency", latency),
			zap.String("clientIP", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("errors", errors.String()),
		}

		switch {
		case statusCode >= 400 && statusCode <= 499:
			zlogger.Warn(c, "GIN", fields...)
		case statusCode >= 500:
			zlogger.Error(c, "GIN", fields...)
		default:
			zlogger.Info(c, "GIN", fields...)
		}
	}
}
