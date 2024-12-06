package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 结构化日志配置
type Config struct {
	ZapConfig        `mapstructure:",squash"`
	LumberjackConfig `mapstructure:",squash"`
	// FieldKeys field 中 key 的白名单
	// 限制 key 的白名单，防止 日志中 key 数量过多
	FieldKeys []string `json:"fieldKeys" yaml:"fieldKeys" mapstructure:"fieldKeys"`
}

// IntoRootLogger 根据配置生成 Root Logger
func (config Config) IntoRootLogger() *zap.Logger {
	// 设置默认值
	config.SetDefaults()

	lumberjackLogger = config.IntoLumberJackLogger(fmt.Sprintf("%s.log", config.Name))

	core := zapcore.NewCore(
		GetEncoder(false),
		zapcore.AddSync(lumberjackLogger),
		config.Level.IntoZapLogLevel(),
	)
	if config.ZapConfig.LogToStdOut {
		core = zapcore.NewTee(
			core,
			zapcore.NewCore(
				GetEncoder(true),
				zapcore.AddSync(os.Stdout),
				config.Level.IntoZapLogLevel(),
			),
		)
	}

	logger := zap.New(core,
		zap.WithCaller(true),
		zap.AddStacktrace(config.TraceLevel.IntoZapLogLevel()),
	).Named(config.Name)

	return logger
}

// SetDefaults 设置默认值
func (config *Config) SetDefaults() {
	if config.Name == "" {
		config.Name = "<untitled>"
	}
	if config.Level == "" {
		config.Level = LogLevelDebug
	}
	if config.TraceLevel == "" {
		config.TraceLevel = LogLevelPanic
	}
}

// GetEncoder 获取日志字段编码器
func GetEncoder(console bool) zapcore.Encoder {
	encodeLevel := zapcore.CapitalLevelEncoder
	if console {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config := zapcore.EncoderConfig{
		NameKey:        "name_",
		TimeKey:        "time_",
		LevelKey:       "level_",
		CallerKey:      "caller_",
		MessageKey:     "msg_",
		StacktraceKey:  "trace_",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if console {
		return zapcore.NewConsoleEncoder(config)
	} else {
		return zapcore.NewJSONEncoder(config)
	}
}
