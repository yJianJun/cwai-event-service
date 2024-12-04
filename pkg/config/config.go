package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

type App struct {
	ConfigFile string `json:"configFile"`
	Port       string `json:"port"`
	Host       string `json:"host"`
}

type ServerConfig struct {
	ElasticSearch ElaticSearch `yaml:"elasticSearch"`
	App           App          `yaml:"app"`
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

var Config *ServerConfig

const ConfigFile = "./conf/config.yaml" // 配置文件

func NewConfig() *ServerConfig {
	Config = &ServerConfig{}
	Config.Init()
	return Config
}

// Init 初始化配置
func (mo *ServerConfig) Init() {
	var configFile string
	// 读取配置文件优先级: 命令行 > 默认值
	flag.StringVar(&configFile, "c", ConfigFile, "配置配置")
	if len(configFile) == 0 {
		// 读取默认配置文件
		fmt.Println("配置文件不存在！")
	}
	// 读取配置文件
	v := viper.New()
	v.SetConfigFile(configFile)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("_", ".")
	viper.SetEnvKeyReplacer(replacer)
	if err := v.ReadInConfig(); err != nil {
		fmt.Errorf("配置解析失败:%s\n", err)
	}
	// 动态监测配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生改变")
		if err := v.Unmarshal(mo); err != nil {
			panic(fmt.Errorf("配置重载失败:%s\n", err))
		}
	})
	if err := v.Unmarshal(mo); err != nil {
		fmt.Errorf("配置重载失败:%s\n", err)
	}
	// 设置配置文件
	mo.App.ConfigFile = configFile
}

func (mo *ServerConfig) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("host", "H", "", "listen host")
	cmd.Flags().StringP("port", "p", "", "listen port")
	viper.BindPFlags(cmd.Flags()) //命令行标志的值设置覆盖值
}

func (mo *ServerConfig) Validate() error {
	//TODO
	return nil
}
