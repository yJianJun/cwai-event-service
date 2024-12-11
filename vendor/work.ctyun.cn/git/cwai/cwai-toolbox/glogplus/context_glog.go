package glogplus

import (
	"context"
	"fmt"

	"github.com/golang/glog"
)

type logPrefixKey struct{}

func SetLogPrefix(ctx context.Context, str string) context.Context {
	return context.WithValue(ctx, logPrefixKey{}, str)
}

func SetLogPrefixf(ctx context.Context, format string, a ...interface{}) context.Context {
	logPrefix := fmt.Sprintf(format, a...)
	return SetLogPrefix(ctx, logPrefix)
}

func AppendLogPrefixf(ctx context.Context, format string, a ...interface{}) context.Context {
	str := fmt.Sprintf(format, a...)
	return AppendLogPrefix(ctx, str)
}

func AppendLogPrefix(ctx context.Context, str string) context.Context {
	logPrefix := GetLogPrefix(ctx) + str
	return SetLogPrefix(ctx, logPrefix)
}

func AppendLogField(ctx context.Context, key string, value interface{}) context.Context {
	if len(key) == 0 {
		return ctx
	}
	if value == nil {
		return AppendLogPrefixf(ctx, "[%s]", key)
	}
	return AppendLogPrefixf(ctx, "[%s:%+v]", key, value)
}

func GetLogPrefix(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(logPrefixKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func Info(ctx context.Context, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	newArgs := append([]interface{}{logPrefix}, args...)
	glog.InfoDepth(1, newArgs...)
}

func Warning(ctx context.Context, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	newArgs := append([]interface{}{logPrefix}, args...)
	glog.InfoDepth(1, newArgs...)
}

func Error(ctx context.Context, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	newArgs := append([]interface{}{logPrefix}, args...)
	glog.InfoDepth(1, newArgs...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	newArgs := append([]interface{}{logPrefix}, args...)
	glog.FatalDepth(1, newArgs...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	s := fmt.Sprintf(format, args...)
	glog.InfoDepth(1, logPrefix, s)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	s := fmt.Sprintf(format, args...)
	glog.WarningDepth(1, logPrefix, s)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	s := fmt.Sprintf(format, args...)
	glog.ErrorDepth(1, logPrefix, s)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logPrefix := GetLogPrefix(ctx)
	s := fmt.Sprintf(format, args...)
	glog.FatalDepth(1, logPrefix, s)
}
