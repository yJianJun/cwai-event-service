package config

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

var EventServerConfig = &ServerConfig{}

type ServerConfig struct {
	App           App          `yaml:"app"`
	ElasticSearch ElaticSearch `yaml:"elasticSearch"`
	LoggerInfo    LoggerInfo   `mapstructure:"logger"`
	AuthInfo      AuthInfo     `mapstructure:"auth"`
}

type App struct {
	Port        string `json:"port"`
	Host        string `json:"host"`
	ShutTimeOut int    `json:"shuttimeout"`
}

type ElaticSearch struct {
	Enable              bool   `yaml:"enable"`
	Url                 string `yaml:"url"`
	Sniff               bool   `yaml:"sniff"`
	HealthCheckInterval int    `yaml:"healthCheckInterval"`
	LogPre              string `yaml:"logPre"`
	Password            string `yaml:"password"`
	Username            string `yaml:"username"`
}

type LoggerInfo struct {
	Name         string `mapstructure:"name"`
	Level        string `mapstructure:"level"`
	TraceLevel   string `mapstructure:"traceLevel"`
	LogToDir     string `mapstructure:"logToDir"`
	MaxSizeInMiB int    `mapstructure:"maxSizeInMiB"`
	MaxAgeInDays int    `mapstructure:"maxAgeInDays"`
}

type AuthInfo struct {
	AuthHost string `mapstructure:"authHost"`
	AuthPath string `mapstructure:"authPath"`
}

func NewConfig() *ServerConfig {
	EventServerConfig.Init()
	if os.Getenv("DEBUG") == "ture" {
		EventServerConfig.LoggerInfo.Level = string(logger.LogLevelDebug)
	}
	return EventServerConfig
}

// Init 初始化默认参数
func (mo *ServerConfig) Init() {
	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", "8081")
	viper.SetDefault("app.shuttimeout", 60)
	viper.SetDefault("logger.name", "cwai-event-service")
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.traceLevel", "error")
	viper.SetDefault("logger.logToDir", ".")
	viper.SetDefault("logger.maxSizeInMiB", "10")
	viper.SetDefault("logger.maxAgeInDays", "30")
	viper.SetDefault("auth.authHost", "https://cwai.ctyun.cn:10004")
	viper.SetDefault("auth.authPath", "/apis/v1/workspace-service/workspace/userInfo")
}

func (mo *ServerConfig) BindFlags(cmd *cobra.Command) error {
	cmd.Flags().StringP("app.host", "H", "", "listen host")
	cmd.Flags().StringP("app.port", "p", "", "listen port")

	err := viper.BindPFlags(cmd.Flags())
	return err
}

// ParseYAML 从YAML文件解析参数
func (mo *ServerConfig) ReadYAML(path string) (err error) {
	viper.SetConfigType("YAML")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(content))

	return err
}

func (mo *ServerConfig) Parse() (err error) {
	return viper.Unmarshal(mo)
}

func (mo *ServerConfig) Validate() error {
	return nil
}
