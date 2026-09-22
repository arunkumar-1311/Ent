package config

import (
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchConfig struct {
	URL       string
	IndexName string
}

func LoadElasticsearchConfig() ElasticsearchConfig {
	url := strings.TrimSpace(os.Getenv("ELASTICSEARCH_URL"))
	if url == "" {
		url = "http://localhost:9200"
	}

	indexName := strings.TrimSpace(os.Getenv("ELASTICSEARCH_SHIPMENTS_INDEX"))
	if indexName == "" {
		indexName = "shipments"
	}

	return ElasticsearchConfig{
		URL:       url,
		IndexName: indexName,
	}
}

func NewElasticsearchClient(
	cfg ElasticsearchConfig,
) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	})
}
