package model

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"fmt"
	"github.com/olivere/elastic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const eventIndexMapping = `{
  "mappings": {
    "_doc": {
      "properties": {
        "bandwidth": {
          "type": "long"
        },
        "data_size": {
          "type": "long"
        },
        "error_code": {
          "type": "long"
        },
        "event_type": {
          "type": "keyword"
        },
        "id": {
          "type": "long"
        },
        "level": {
          "type": "keyword"
        },
        "local_guid": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_max_word"
        },
        "remote_guid": {
          "type": "text",
          "analyzer": "ik_max_word",
          "search_analyzer": "ik_max_word"
        },
        "time_duration": {
          "type": "long"
        },
        "timestamp": {
          "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
        }
      }
    }
  }
}`

var DB *gorm.DB
var ESclient *elastic.Client // ES客户端

func ConnectDatabase(conf *config.ServerConfig) error {
	dsn := generateDSN(conf.MysqlConfig)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := migrateDatabase(database); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = database
	return nil
}

// 生成DSN字符串
func generateDSN(mysqlConfig config.MysqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Address, mysqlConfig.Db)
}

// 迁移数据库
func migrateDatabase(database *gorm.DB) error {
	return database.AutoMigrate(&Event{})
}

func InitElasticSearch(conf *config.ServerConfig) error {
	elasticConfig := conf.ElasticConfig
	if !elasticConfig.Enable {
		return nil
	}

	fmt.Printf("elasticConfig: %v\n", elasticConfig)
	client, err := createElasticClient(elasticConfig)
	if err != nil {
		return fmt.Errorf("failed to init elasticsearch: %w", err)
	}

	ESclient = client
	if err := createIndexMapping(ESclient); err != nil {
		return fmt.Errorf("failed to create index mapping: %w", err)
	}

	return nil
}

func createElasticClient(config config.ElasticConfig) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Url),
		elastic.SetSniff(config.Sniff),
		elastic.SetHealthcheckInterval(config.HealthCheckInterval),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func createIndexMapping(client *elastic.Client) error {
	return checkAndCreateIndex(client, "events", eventIndexMapping)
}

func checkAndCreateIndex(client *elastic.Client, indexName string, mapping string) error {
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		_, err := client.CreateIndex(indexName).Body(mapping).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
