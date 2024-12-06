package logger

import "go.uber.org/zap/zapcore"

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelPanic LogLevel = "panic"
	LogLevelFatal LogLevel = "fatal"
)

func (level LogLevel) IntoZapLogLevel() zapcore.Level {
	switch level {
	case LogLevelDebug:
		return zapcore.DebugLevel
	case LogLevelInfo:
		return zapcore.InfoLevel
	case LogLevelWarn:
		return zapcore.WarnLevel
	case LogLevelError:
		return zapcore.ErrorLevel
	case LogLevelPanic:
		return zapcore.PanicLevel
	case LogLevelFatal:
		return zapcore.FatalLevel
	default:
		panic("unsupported log level")
	}
}
