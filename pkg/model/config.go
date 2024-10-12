package model

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// type LoggerInfo struct {
// 	Name         string
// 	Level        string
// 	TraceLevel   string
// 	LogToDir     string
// 	MaxSizeInMiB int
// 	MaxAgeInDays int
// }

type CCAEServer struct {
	Protocol     string `mapstructure:"protocol"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	UserName     string `mapstructure:"userName"`
	UserPassword string `mapstructure:"userPassword"`
	TokenTimeOut int    `mapstructure:"tokenTimeOut"`
}

type CCAEAPIs struct {
	TokenUrl     string `mapstructure:"tokenUrl"`
	QureyTopoUrl string `mapstructure:"qureyTopoUrl"`
}

type ServerConfig struct {
	Protocol      string        `mapstructure:"protocol"`
	Host          string        `mapstructure:"host"`
	Port          string        `mapstructure:"port"`
	ShutTimeOut   int           `mapstructure:"shutTimeOut"`
	CCAEServer    CCAEServer    `mapstructure:"ccaeServer"`
	CCAEAPIs      CCAEAPIs      `mapstructure:"ccaeAPIs"`
	MysqlConfig   MysqlConfig   `mapstructure:"mysqlConfig"`
	ElasticConfig ElasticConfig `mapstructure:"elasticConfig"`
	//	LoggerInfo  LoggerInfo `mapstructure:"logger"`
}

type MysqlConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	Db       string `mapstructure:"db"`
}

type ElasticConfig struct {
	Enable              bool          `mapstructure:"enable"`
	Url                 string        `mapstructure:"url"`
	Sniff               bool          `mapstructure:"sniff"`
	HealthCheckInterval time.Duration `mapstructure:"healthCheckInterval"`
	LogPre              string        `mapstructure:"logPre"`
}

func NewConfig() *ServerConfig {
	config := &ServerConfig{}
	config.Init()
	return config
}

// Init 初始化默认参数
func (mo *ServerConfig) Init() {
	viper.SetDefault("protocol", "http")
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "8081")
	viper.SetDefault("shutTimeOut", "5")
	viper.SetDefault("ccaeServer.protocol", "https")
	viper.SetDefault("ccaeServer.host", "ccae.manager.ctyun.cn")
	viper.SetDefault("ccaeServer.port", "26335")
	viper.SetDefault("ccaeServer.userName", "lndxsy502cs")
	viper.SetDefault("ccaeServer.userPassword", "Lnsydx502cs@2024")
	viper.SetDefault("ccaeServer.tokenTimeOut", "1800")
	viper.SetDefault("ccaeAPIs.tokenUrl", "/rest/plat/smapp/v1/oauth/token")
	viper.SetDefault("ccaeAPIs.qureyTopoUrl", "/api/monitor/v1/query-topo")
	viper.SetDefault("mysqlConfig.user", "roor")
	viper.SetDefault("mysqlConfig.password", "930927")
	viper.SetDefault("mysqlConfig.address", "127.0.0.1:3306")
	viper.SetDefault("mysqlConfig.db", "ctccl")
}

func (mo *ServerConfig) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("host", "H", "", "listen host")
	cmd.Flags().StringP("port", "p", "", "listen port")

	viper.BindPFlags(cmd.Flags()) //命令行标志的值设置覆盖值
}

// ParseYAML 从YAML文件解析参数
func (mo *ServerConfig) ReadYAML(path string) (err error) {
	viper.SetConfigType("YAML")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	viper.ReadConfig(bytes.NewBuffer(content))
	return nil
}

func (mo *ServerConfig) Parse() (err error) {
	return viper.Unmarshal(mo)
}

func (mo *ServerConfig) Validate() error {
	//TODO
	return nil
}
