package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"work.ctyun.cn/git/cwai/cwai-toolbox/utils"
)

var logger *zap.Logger
var wrappedLogger *zap.Logger
var lumberjackLogger *lumberjack.Logger
var fieldKeySet map[string]bool

type loggerKey struct{}

const ginLoggerKey = "ginLogger"
const traceIDKey = "TraceID"
const funcKey = "Func"
const tagKey = "Tag"

func init() {
	logger = zap.New(zapcore.NewCore(GetEncoder(true), os.Stdout, zapcore.InfoLevel))
	wrappedLogger = logger.WithOptions(zap.AddCallerSkip(1))
}

func Init(config *Config) {
	logger = config.IntoRootLogger()

	wrappedLogger = logger.WithOptions(zap.AddCallerSkip(1))
	logger.Info("logger initialized with config",
		Stringify("config", config),
	)

	fieldKeySet = map[string]bool{}
	for _, v := range config.FieldKeys {
		fieldKeySet[v] = true
	}
}

func GetLumberjackLogger() *lumberjack.Logger {
	return lumberjackLogger
}

// Stringify converts anything into a string, in JSON form.
// this is primarily used to prevent field booming in elasticsearch.
func Stringify(name string, any interface{}) zap.Field {
	json, err := jsoniter.MarshalToString(any)
	if err != nil {
		return zap.NamedError(name, err)
	}
	return zap.String(name, json)
}

func appendTraceId(ctx context.Context, format string) string {
	var builder strings.Builder

	traceID := ctx.Value(traceIDKey)
	if traceID != nil {
		builder.WriteString(traceID.(string))
		builder.WriteString(" ")
	}

	funcID := ctx.Value(funcKey)
	if funcID != nil {
		builder.WriteString(funcID.(string))
		builder.WriteString(" ")
	}

	tagID := ctx.Value(tagKey)
	if tagID != nil {
		builder.WriteString(tagID.(string))
		builder.WriteString(" ")
	}

	builder.WriteString(format)
	return builder.String()
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	getLogger(ctx).Debug(appendTraceId(ctx, msg), filterFields(ctx, fields)...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(appendTraceId(ctx, format), args...)
	getLogger(ctx).Debug(msg)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	getLogger(ctx).Info(appendTraceId(ctx, msg), filterFields(ctx, fields)...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(appendTraceId(ctx, format), args...)
	getLogger(ctx).Info(msg)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	getLogger(ctx).Warn(appendTraceId(ctx, msg), filterFields(ctx, fields)...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(appendTraceId(ctx, format), args...)
	getLogger(ctx).Warn(msg)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	getLogger(ctx).Error(appendTraceId(ctx, msg), filterFields(ctx, fields)...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(appendTraceId(ctx, format), args...)
	getLogger(ctx).Error(msg)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	getLogger(ctx).Panic(appendTraceId(ctx, msg), filterFields(ctx, fields)...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(appendTraceId(ctx, format), args...)
	getLogger(ctx).Panic(msg)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	getLogger(ctx).Fatal(appendTraceId(ctx, msg), filterFields(ctx, fields)...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(appendTraceId(ctx, format), args...)
	getLogger(ctx).Fatal(msg)
}

func L() *zap.Logger {
	return logger
}

func ContextWithLogFields(ctx context.Context, fields ...zap.Field) context.Context {
	logger := getLogger(ctx)
	newLogger := logger.With(filterFields(ctx, fields)...)
	ctx = context.WithValue(ctx, loggerKey{}, newLogger)
	return context.WithValue(ctx, traceIDKey, buildTraceId())
}

func InjectLogFieldsForContext(ctx context.Context, tagId, funcID string) context.Context {
	ctx = context.WithValue(ctx, funcKey, funcID)
	return context.WithValue(ctx, tagKey, tagId)
}

func SetLogFieldsForGinContext(ctx *gin.Context, fields ...zap.Field) {
	logger := getLogger(ctx)
	newLogger := logger.With(filterFields(ctx, fields)...)
	ctx.Set(ginLoggerKey, newLogger)
	ctx.Set(traceIDKey, buildTraceId())
}
func InjectLogFieldsForGinContext(ctx *gin.Context, tagId, funcID string) {
	ctx.Set(funcKey, funcID)
	ctx.Set(tagKey, tagId)
}

func getLogger(ctx context.Context) *zap.Logger {
	v := ctx.Value(loggerKey{})
	if v != nil {
		return v.(*zap.Logger)
	}
	v = ctx.Value(ginLoggerKey)
	if v != nil {
		return v.(*zap.Logger)
	}
	// can not find logger in context, using global logger instead
	return wrappedLogger
}

// filter fields which is not whitelisted to avoid field blasing
func filterFields(ctx context.Context, fields []zap.Field) []zap.Field {
	if len(fieldKeySet) == 0 {
		return fields
	}
	if len(fields) == 0 {
		return []zap.Field{}
	}
	filteredFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		if _, find := fieldKeySet[f.Key]; find {
			filteredFields = append(filteredFields, f)
		} else {
			Warnf(ctx, "ignore field(key:%s) in context logger: not whitelisted", f.Key)
		}
	}
	return filteredFields
}

func Flush() {
	logger.Sync()
	wrappedLogger.Sync()
}

func buildTraceId() string {
	return utils.FormatLogNow() + uuid.New().String()[0:8]
}
