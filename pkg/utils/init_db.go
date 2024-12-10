package utils

import (
	"context"
	"errors"

	"github.com/elastic/go-elasticsearch/v8"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/config"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

var ESclient *elasticsearch.TypedClient // ES客户端

func InitElasticSearch() error {
	elasticConfig := config.EventServerConfig.ElasticSearch
	if !elasticConfig.Enable {
		return errors.New("elasticsearch disabled")
	}

	logger.Infof(context.TODO(), "elasticConfig: %v\n", elasticConfig)
	client, err := createElasticClient(elasticConfig)
	if err != nil {
		return err
	}
	ESclient = client

	return nil
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
