package middleware

import (
	zlogger "work.ctyun.cn/git/cwai/cwai-toolbox/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type glogPanicWriter struct{}

func (pw glogPanicWriter) Write(p []byte) (int, error) {
	glog.Error(string(p))
	return len(p), nil
}

var glogWriter glogPanicWriter

func RecoveryWithGlog() gin.HandlerFunc {
	return gin.RecoveryWithWriter(glogWriter)
}

type zapPanicWriter struct {
	ctx *gin.Context
}

func (pw zapPanicWriter) Write(p []byte) (int, error) {
	zlogger.Error(pw.ctx, string(p))
	return len(p), nil
}

func RecoveryWithZap() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gin.RecoveryWithWriter(zapPanicWriter{ctx: ctx})(ctx)
	}
}
