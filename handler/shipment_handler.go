package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"shipment-service/model"
	"shipment-service/service"
)

type ShipmentHandler struct {
	service *service.ShipmentService
}

func NewShipmentHandler(
	service *service.ShipmentService,
) *ShipmentHandler {

	return &ShipmentHandler{
		service: service,
	}
}

func (h *ShipmentHandler) Create(c echo.Context) error {

	var req model.CreateShipmentRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid request",
			},
		)
	}

	shipment, err := h.service.CreateShipment(
		c.Request().Context(),
		req,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(http.StatusCreated, shipment)
}

func (h *ShipmentHandler) GetByID(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid shipment ID",
			},
		)
	}

	shipment, err := h.service.GetShipment(
		c.Request().Context(),
		id,
	)

	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			map[string]string{
				"error": "shipment not found",
			},
		)
	}

	return c.JSON(http.StatusOK, shipment)
}


func (h *ShipmentHandler) GetAll(c echo.Context) error {

	shipments, err := h.service.GetShipments(
		c.Request().Context(),
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, shipments)
}

func (h *ShipmentHandler) Update(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid shipment ID",
			},
		)
	}

	var req model.UpdateShipmentRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid request",
			},
		)
	}

	shipment, err := h.service.UpdateShipment(
		c.Request().Context(),
		id,
		req,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, shipment)
}

func (h *ShipmentHandler) Delete(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid shipment ID",
			},
		)
	}

	err = h.service.DeleteShipment(
		c.Request().Context(),
		id,
	)

	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			map[string]string{
				"error": "shipment not found",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "shipment deleted successfully",
		},
	)
}