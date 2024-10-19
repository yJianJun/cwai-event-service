package config

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

type server struct {
	Port            string `yaml:"port"`
	Protocol        string `yaml:"protocol"`
	Host            string `yaml:"host"`
	UserName        string `yaml:"username"`
	UserPassword    string `yaml:"password"`
	TokenTimeOut    int64  `yaml:"tokenTimeOut"`
	ShutdownTimeout int64  `yaml:"shutdownTimeout"`
}

type api struct {
	TokenUrl     string `yaml:"tokenUrl"`
	QureyTopoUrl string `yaml:"qureyTopoUrl"`
}

type CCAE struct {
	Server server `yaml:"server"`
	Api    api    `yaml:"api"`
}

type App struct {
	ConfigFile string `json:"configFile"`
	Port       string `json:"port"`
	Host       string `json:"host"`
}

type ServerConfig struct {
	CCAE          CCAE         `yaml:"ccae"`
	Mysql         Mysql        `yaml:"mysql"`
	ElasticSearch ElaticSearch `yaml:"elasticSearch"`
	Log           log          `yaml:"log"`
	App           App          `yaml:"app"`
}

type Mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Db       string `yaml:"db"`
}

type ElaticSearch struct {
	Enable              bool   `yaml:"enable"`
	Url                 string `yaml:"url"`
	Sniff               bool   `yaml:"sniff"`
	HealthCheckInterval int    `yaml:"healthCheckInterval"`
	LogPre              string `yaml:"logPre"`
}

type log struct {
	Dir string `yaml:"dir"`
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
