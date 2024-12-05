package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	KMLEnvVarPodName    = "MY_POD_NAME"
	KMLEnvVarPodVersion = "K8S_POD_VERSION"
)

var (
	ErrKMLEnvVarNotExist = errors.New("环境变量不存在")
)

type PodInfo struct {
	PodName    string `json:"pod_name"`
	PodVersion string `json:"pod_version"`
}

// InjectIntoZapLogger 将信息注入 zap logger
func (info PodInfo) InjectIntoZapLogger(logger *zap.Logger) *zap.Logger {
	return logger.With(
		zap.String("pod_name", info.PodName),
		zap.String("pod_version", info.PodVersion),
	)
}
