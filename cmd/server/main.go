package main

import (
	"context"
	"log"

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

	// Repository
	shipmentRepository := repository.NewShipmentRepository(client)

	// Service
	shipmentService := service.NewShipmentService(
		shipmentRepository,
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

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
