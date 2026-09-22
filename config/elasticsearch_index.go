package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
)

func EnsureShipmentIndex(
	ctx context.Context,
	client *elasticsearch.Client,
	indexName string,
) error {
	existsResponse, err := client.Indices.Get(
		[]string{indexName},
		client.Indices.Get.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("check shipment index: %w", err)
	}
	defer existsResponse.Body.Close()

	if existsResponse.StatusCode == 200 {
		return nil
	}

	if existsResponse.StatusCode != 404 {
		body, readErr := io.ReadAll(existsResponse.Body)
		if readErr != nil {
			return fmt.Errorf(
				"check shipment index returned status %d: %w",
				existsResponse.StatusCode,
				readErr,
			)
		}

		return fmt.Errorf(
			"check shipment index returned status %d: %s",
			existsResponse.StatusCode,
			string(body),
		)
	}

	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":              map[string]string{"type": "integer"},
				"tracking_number": map[string]string{"type": "keyword"},
				"sender_name":     map[string]string{"type": "text"},
				"receiver_name":   map[string]string{"type": "text"},
				"origin":          map[string]string{"type": "text"},
				"destination":     map[string]string{"type": "text"},
				"status":          map[string]string{"type": "keyword"},
				"weight":          map[string]string{"type": "double"},
			},
		},
	}

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("encode shipment index mapping: %w", err)
	}

	createResponse, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithContext(ctx),
		client.Indices.Create.WithBody(bytes.NewReader(mappingJSON)),
	)
	if err != nil {
		return fmt.Errorf("create shipment index: %w", err)
	}
	defer createResponse.Body.Close()

	if createResponse.IsError() {
		body, readErr := io.ReadAll(createResponse.Body)
		if readErr != nil {
			return fmt.Errorf(
				"create shipment index returned status %d: %w",
				createResponse.StatusCode,
				readErr,
			)
		}

		return fmt.Errorf(
			"create shipment index returned status %d: %s",
			createResponse.StatusCode,
			string(body),
		)
	}

	return nil
}
