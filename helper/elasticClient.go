package helper

import (
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v9"
)

func NewElasticClient() (*elasticsearch.Client, error) {
	elasticHost := os.Getenv("ELASTICHOST")
	if elasticHost == "" {
		elasticHost = "localhost"
	}

	address := fmt.Sprintf("http://%s:9200", elasticHost)
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := esClient.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch returned an error: %s", res.Status())
	}

	return esClient, nil
}
