package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type logger interface {
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

// glog implementation
type glogger struct {
}

func (glogger) Info(args ...interface{}) {
	glog.Info(args...)
}

func (glogger) Warning(args ...interface{}) {
	glog.Warning(args...)
}

func (glogger) Error(args ...interface{}) {
	glog.Error(args...)
}

func setupLogging(duration time.Duration) {
	go func() {
		for range time.Tick(duration) {
			glog.Flush()
		}
	}()
}

// Logger middleware
// Example:
//
//	router := gin.New()
//	router.Use(ginglog.Logger(3 * time.Second))
func Logger(duration time.Duration) gin.HandlerFunc {
	setupLogging(duration)
	return LoggerHandler(glogger{})
}
