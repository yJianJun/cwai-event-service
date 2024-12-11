package logger_test

import (
	"context"
	"io"
	"testing"
	"time"

	"go.uber.org/zap"
	log "work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

var (
	ctx context.Context
)

func init() {
	log.Init(&log.Config{
		ZapConfig: log.ZapConfig{
			Name:        "solaria-test",
			Level:       log.LogLevelInfo,
			TraceLevel:  log.LogLevelError,
			LogToStdOut: true,
		},
		LumberjackConfig: log.LumberjackConfig{
			LogToDir:     "/tmp",
			MaxSizeInMiB: 10,
			MaxAgeInDays: 28,
		},
		FieldKeys: []string{"reason", "ID", "who", "what", "error", "config"},
	})
	ctx = context.Background()
}

func TestEncapsulatedLogger(t *testing.T) {
	log.Debug(ctx, "i am invisible",
		zap.String("reason", "min log level"))
	log.Info(ctx, "encapsulated test",
		zap.String("hello", "world"),
		zap.Int("ID", 114514))
	log.Warn(ctx, "hello",
		zap.String("who", "you"),
		zap.Error(io.ErrUnexpectedEOF))
	log.Error(ctx, "shit",
		zap.String("what", "this is terrible"),
		zap.Error(io.ErrClosedPipe),
		zap.String("unknowKey", "should not display"),
	)
}

func TestColorfulLog(t *testing.T) {
	var (
		green  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
		white  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
		yellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
		red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
		reset  = string([]byte{27, 91, 48, 109})
	)
	log.Infof(ctx, "look! %s", "this is common log")
	log.Infof(ctx, "look! %s%s%s", green, "this is green log", reset)
	log.Infof(ctx, "look! %s%s%s", white, "this is white log", reset)
	log.Infof(ctx, "look! %s%s%s", yellow, "this is yellow log", reset)
	log.Errorf(ctx, "look! %s%s%s", red, "this is red log", reset)
}

func TestLoggerInContext(t *testing.T) {
	context := log.ContextWithLogFields(ctx, zap.String("reason", "min log level"))
	log.Debug(context, "i am invisible")
	log.Debugf(context, "i am invisible %s", "too")

	context = log.ContextWithLogFields(ctx,
		zap.String("hello", "world"),
		zap.Int("ID", 114514))
	log.Info(context, "encapsulated test")
	log.Infof(context, "encapsulated %s test", "format")

	context = log.ContextWithLogFields(ctx,
		zap.String("who", "you"))
	log.Warn(context, "hello", zap.Error(io.ErrUnexpectedEOF))
	log.Warnf(context, "hello %d", 110)

	context = log.ContextWithLogFields(ctx,
		zap.String("what", "this is terrible"))
	log.Error(context, "shit", zap.Error(io.ErrClosedPipe))
	log.Errorf(context, "shit %v", time.Now())
}

func TestNestedOutput(t *testing.T) {
	deepDarkStruct := struct {
		DeeperDarkerField map[string]map[string]string
	}{
		DeeperDarkerField: map[string]map[string]string{
			"hell": {
				"this": "is a nested hell!",
				"and":  "will make elasticsearch boom",
			},
		},
	}
	log.Error(ctx, "this is a nested output", zap.Any("config", deepDarkStruct))
	log.Error(ctx, "this is a de-nested output", log.Stringify("config", deepDarkStruct))
}
