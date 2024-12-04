package util

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

var ESclient *elasticsearch.TypedClient // ES客户端

func InitElasticSearch(conf *config.ServerConfig) {
	elasticConfig := conf.ElasticSearch
	if !elasticConfig.Enable {
		return
	}

	fmt.Printf("elasticConfig: %v\n", elasticConfig)
	client, err := createElasticClient(elasticConfig)
	if err != nil {
		common.Error((map[string]interface{}{
			"message": "failed to init elasticsearch",
			"err":     err,
		}))
	}
	ESclient = client
}

func createElasticClient(esConfig config.ElaticSearch) (*elasticsearch.TypedClient, error) {
	elasticConfig := elasticsearch.Config{
		Addresses: []string{esConfig.Url},
		Username:  esConfig.Username,
		Password:  esConfig.Password,
	}
	client, err := elasticsearch.NewTypedClient(elasticConfig)
	if err != nil {
		return nil, err
	}
	return client, err
}
