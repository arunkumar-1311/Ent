package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"

	"shipment-service/ent"
	"shipment-service/model"
)

type ShipmentSearchRepository struct {
	client    *elasticsearch.Client
	indexName string
}

func NewShipmentSearchRepository(
	client *elasticsearch.Client,
	indexName string,
) *ShipmentSearchRepository {
	return &ShipmentSearchRepository{
		client:    client,
		indexName: indexName,
	}
}

func (r *ShipmentSearchRepository) Index(
	ctx context.Context,
	shipment *ent.Shipment,
) error {
	document := shipmentDocument(shipment)

	body, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("encode shipment document: %w", err)
	}

	response, err := r.client.Index(
		r.indexName,
		bytes.NewReader(body),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(strconv.Itoa(shipment.ID)),
		r.client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return fmt.Errorf("index shipment %d: %w", shipment.ID, err)
	}
	defer response.Body.Close()

	if response.IsError() {
		responseBody, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return fmt.Errorf(
				"index shipment %d returned status %d: %w",
				shipment.ID,
				response.StatusCode,
				readErr,
			)
		}

		return fmt.Errorf(
			"index shipment %d returned status %d: %s",
			shipment.ID,
			response.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (r *ShipmentSearchRepository) Delete(
	ctx context.Context,
	shipmentID int,
) error {
	response, err := r.client.Delete(
		r.indexName,
		strconv.Itoa(shipmentID),
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("wait_for"),
	)
	if err != nil {
		return fmt.Errorf("delete shipment %d from index: %w", shipmentID, err)
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return nil
	}

	if response.IsError() {
		responseBody, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return fmt.Errorf(
				"delete shipment %d returned status %d: %w",
				shipmentID,
				response.StatusCode,
				readErr,
			)
		}

		return fmt.Errorf(
			"delete shipment %d returned status %d: %s",
			shipmentID,
			response.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (r *ShipmentSearchRepository) Search(
	ctx context.Context,
	query string,
) ([]model.ShipmentDocument, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
				"fields": []string{
					"tracking_number",
					"sender_name",
					"receiver_name",
					"origin",
					"destination",
					"status",
				},
			},
		},
	}

	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("encode shipment search query: %w", err)
	}

	response, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.indexName),
		r.client.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("search shipments: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		responseBody, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return nil, fmt.Errorf(
				"search shipments returned status %d: %w",
				response.StatusCode,
				readErr,
			)
		}

		return nil, fmt.Errorf(
			"search shipments returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var searchResponse struct {
		Hits struct {
			Hits []struct {
				Source model.ShipmentDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("decode shipment search response: %w", err)
	}

	shipments := make([]model.ShipmentDocument, 0, len(searchResponse.Hits.Hits))

	for _, hit := range searchResponse.Hits.Hits {
		shipments = append(shipments, hit.Source)
	}

	return shipments, nil
}

func shipmentDocument(shipment *ent.Shipment) model.ShipmentDocument {
	return model.ShipmentDocument{
		ID:             shipment.ID,
		TrackingNumber: shipment.TrackingNumber,
		SenderName:     shipment.SenderName,
		ReceiverName:   shipment.ReceiverName,
		Origin:         shipment.Origin,
		Destination:    shipment.Destination,
		Status:         shipment.Status,
		Weight:         shipment.Weight,
	}
}
