package model

import (
	"context"
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
		err = CreateIndexMapping(ESclient)
		if err != nil {
			return fmt.Errorf("failed to create index mapping: %w", err)
		}
		return nil
	}
	return nil
}

func CreateIndexMapping(client *elastic.Client) error {
	exists, err := client.IndexExists("events").Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		// 配置映射
		mapping := `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "_doc": {
      "properties": {
        "localguid": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_max_word"
        },
        "remoteguid": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_max_word"
        },
        "errcode": {
          "type": "long"
        },
        "timeduration": {
          "type": "long"
        },
        "datasize": {
          "type": "long"
        },
        "bandwidth": {
          "type": "integer"
        },
        "time": {
          "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
        },
        "level": {
          "type": "keyword"
        },
        "event_type": {
          "type": "keyword",
          "index": false
        }
      }
    }
  }
}`
		//注意：增加的写法-创建索引配置映射
		_, err := client.CreateIndex("events").Body(mapping).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
