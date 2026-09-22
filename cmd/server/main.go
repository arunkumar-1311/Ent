package main

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"shipment-service/config"
	"shipment-service/handler"
	"shipment-service/repository"
	"shipment-service/service"
)

func main() {

	// Database
	client := config.ConnectDatabase()
	defer client.Close()

	// Migration
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal("failed to create schema:", err)
	}

	// Elasticsearch
	elasticsearchConfig := config.LoadElasticsearchConfig()

	elasticsearchClient, err := config.NewElasticsearchClient(
		elasticsearchConfig,
	)
	if err != nil {
		log.Fatal("failed to create Elasticsearch client:", err)
	}

	elasticsearchContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := config.EnsureShipmentIndex(
		elasticsearchContext,
		elasticsearchClient,
		elasticsearchConfig.IndexName,
	); err != nil {
		log.Fatal("failed to initialize Elasticsearch index:", err)
	}

	// Repository
	shipmentRepository := repository.NewShipmentRepository(client)
	shipmentSearchRepository := repository.NewShipmentSearchRepository(
		elasticsearchClient,
		elasticsearchConfig.IndexName,
	)

	// Service
	shipmentService := service.NewShipmentService(
		shipmentRepository,
		shipmentSearchRepository,
	)

	// Handler
	shipmentHandler := handler.NewShipmentHandler(
		shipmentService,
	)

	// Echo
	e := echo.New()

	// Routes
	e.POST("/shipments", shipmentHandler.Create)
	e.GET("/shipments", shipmentHandler.GetAll)
	e.GET("/shipments/:id", shipmentHandler.GetByID)
	e.PUT("/shipments/:id", shipmentHandler.Update)
	e.DELETE("/shipments/:id", shipmentHandler.Delete)

	e.GET("/shipments/search", shipmentHandler.Search)
    e.POST("/shipments/reindex", shipmentHandler.Reindex)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
