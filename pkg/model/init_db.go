package model

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"strings"
)

const event_mapping = `{
  "mappings": {
    "properties": {
      "spec_version": {
        "type": "keyword"
      },
      "id": {
        "type": "keyword"
      },
      "source": {
        "type": "keyword"
      },
      "ctyun_region": {
        "type": "keyword"
      },
      "type": {
        "type": "keyword"
      },
      "data_content_type": {
        "type": "keyword"
      },
      "subject": {
        "type": "keyword"
      },
      "time": {
        "type": "date",
        "format": "strict_date_optional_time||yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      },
      "data": {
        "properties": {
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
            "type": "text"
          },
          "account_id": {
            "type": "keyword"
          },
          "user_id": {
            "type": "keyword"
          },
          "compute_type": {
            "type": "keyword"
          },
          "node_ip": {
            "type": "ip"
          },
          "node_name": {
            "type": "keyword"
          },
          "pod_namespace": {
            "type": "keyword"
          },
          "pod_ip": {
            "type": "ip"
          },
          "pod_name": {
            "type": "keyword"
          },
          "region_id": {
            "type": "keyword"
          },
          "resource_group_id": {
            "type": "keyword"
          },
          "resource_group_name": {
            "type": "keyword"
          },
          "level": {
            "type": "keyword"
          },
          "status": {
            "type": "keyword"
          },
          "event_massage": {
            "type": "text"
          },
          "local_guid": {
            "type": "keyword"
          },
          "remote_guid": {
            "type": "keyword"
          },
          "err_code": {
            "type": "keyword"
          },
          "err_message": {
            "type": "text"
          },
          "status_message": {
            "type": "text"
          }
        }
      }
    }
  }
}`

var ESclient *elasticsearch.TypedClient // ES客户端

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

func createElasticClient(esConfig config.ElaticSearch) (*elasticsearch.TypedClient, error) {
	elasticConfig := elasticsearch.Config{Addresses: []string{esConfig.Url}}
	client, err := elasticsearch.NewTypedClient(elasticConfig)
	if err != nil {
		return nil, err
	}
	return client, err
}

func createIndexMapping(client *elasticsearch.TypedClient) error {
	// 创建索引和映射时捕获错误
	if err := checkAndCreateIndex(client, "events", event_mapping); err != nil {
		return fmt.Errorf("failed to create compute task event mapping: %w", err)
	}
	return nil
}

func checkAndCreateIndex(client *elasticsearch.TypedClient, indexName string, mapping string) error {
	if exists, err := client.Indices.Exists(indexName).IsSuccess(context.Background()); !exists && err == nil {
		res, err := client.Indices.Create(indexName).Raw(strings.NewReader(mapping)).Do(context.Background())
		if err != nil {
			return err
		}
		if !res.Acknowledged && res.Index != indexName {
			common.Error(map[string]interface{}{"err": res}, "unexpected error during index creation")
		}
	} else if err != nil {
		return err
	}
	return nil
}
