package logger

import (
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LumberjackConfig 封装的 Lumberjack 的配置。
// Lumberjack 是一个负责自动将日志按大小和日期归档的第三方库。
type LumberjackConfig struct {
	// LogToDir 日志文件的输出路径。
	LogToDir string `json:"logToDir" yaml:"logToDir" mapstructure:"logToDir"`
	// MaxSizeInMiB 单个日志文件最大的大小，超出后会被 rotate 成另一个文件。
	MaxSizeInMiB int `json:"maxSizeInMiB" yaml:"maxSizeInMiB" mapstructure:"maxSizeInMiB"`
	// MaxAgeInDays 日志文件最长的保存时限，过期的日志会被自动删除。
	MaxAgeInDays int `json:"maxAgeInDays" yaml:"maxAgeInDays" mapstructure:"maxAgeInDays"`
}

// IntoLumberJackLogger 根据 LumberjackConfig 构建 Lumberjack 对象。
func (cfg LumberjackConfig) IntoLumberJackLogger(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename: path.Join(cfg.LogToDir, filename),
		MaxSize:  cfg.MaxSizeInMiB,
		MaxAge:   cfg.MaxAgeInDays,
	}
}
