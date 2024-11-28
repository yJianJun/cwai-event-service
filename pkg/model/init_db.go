package model

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

const eventIndexMapping = `{
  "mappings": {
    "_doc": {
      "_size": {
        "enabled": true
      },
      "properties": {
        "specversion": {
          "type": "keyword"
        },
        "id": {
          "type": "keyword"
        },
        "source": {
          "type": "keyword"
        },
        "ctyunregion": {
          "type": "keyword"
        },
        "type": {
          "type": "keyword"
        },
        "datacontenttype": {
          "type": "keyword"
        },
        "subject": {
          "type": "keyword"
        },
        "time": {
          "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
        },
        "task_id": {
          "type": "keyword"
        },
        "task_record_id": {
          "type": "keyword"
        },
        "task_name": {
          "type": "keyword"
        },
        "task_detail": {
          "type": "keyword"
        },
        "level": {
          "type": "keyword"
        },
        "account_id": {
          "type": "keyword"
        },
        "user_id": {
          "type": "keyword"
        },
        "region_id": {
          "type": "keyword"
        },
        "resource_group_id": {
          "type": "keyword"
        },
        "resource_group_name": {
          "type": "text",
          "analyzer": "ik_max_word"
        },
        "data": {
          "properties": {
            "time": {
              "type": "date",
              "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
            },
            "status": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "status_message": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "compute_type": {
              "type": "keyword"
            },
            "node_ip": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "node_name": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "pod_namespace": {
              "type": "keyword"
            },
            "pod_ip": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "pod_name": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "compute_detail": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "local_guid": { "type": "keyword"},
            "remote_guid": { "type": "keyword"},
            "errcode": { "type": "keyword"},
            "err_message": {
              "type": "text",
              "analyzer": "ik_max_word"
            },
            "content": {
              "type": "text",
              "analyzer": "ik_max_word"
            }
          }
        }
      }
    }
  }
}`

var ESclient *elastic.Client // ES客户端

// 生成DSN字符串
func generateDSN(mysqlConfig config.Mysql) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Address, mysqlConfig.Db)
}

func InitElasticSearch(conf *config.ServerConfig) error {
	elasticConfig := conf.ElasticSearch
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

func createElasticClient(config config.ElaticSearch) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Url),
		elastic.SetSniff(config.Sniff),
		elastic.SetHealthcheckInterval(time.Duration(config.HealthCheckInterval)),
	)
	service := elastic.NewXPackSecurityChangePasswordService(client)
	_, err = service.Password(config.Password).Username("elastic").Do(context.Background())
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
