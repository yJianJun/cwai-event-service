package model

import (
	"fmt"
	"github.com/olivere/elastic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ESclient *elastic.Client // ES客户端

func ConnectDatabase(conf *ServerConfig) error {
	dsn := conf.MysqlConfig.User + ":" + conf.MysqlConfig.Password + "@tcp(" + conf.MysqlConfig.Address + ")/" + conf.MysqlConfig.Db + "?charset=utf8&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := database.AutoMigrate(&Event{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	DB = database
	return nil
}

func InitES(conf *ServerConfig) error {
	// 配置
	elasticConfig := conf.ElasticConfig
	if elasticConfig.Enable {
		fmt.Printf("elasticConfig: %v\n", elasticConfig)
		// 创建客户端
		client, err := elastic.NewClient(
			elastic.SetURL(elasticConfig.Url),
			elastic.SetSniff(elasticConfig.Sniff),
			elastic.SetHealthcheckInterval(elasticConfig.HealthCheckInterval),
		)
		if err != nil {
			return fmt.Errorf("failed to init elasticsearch: %w", err)
		}
		ESclient = client
		return nil
	}
	return nil
}
