package service

import (
	"context"

	"shipment-service/ent"
	"shipment-service/model"
	"shipment-service/repository"
)

type ShipmentService struct {
	repository *repository.ShipmentRepository
}

func NewShipmentService(
	repository *repository.ShipmentRepository,
) *ShipmentService {

	return &ShipmentService{
		repository: repository,
	}
}

func (s *ShipmentService) CreateShipment(
	ctx context.Context,
	req model.CreateShipmentRequest,
) (*ent.Shipment, error) {

	return s.repository.Create(ctx, req)
}

func (s *ShipmentService) GetShipment(
	ctx context.Context,
	id int,
) (*ent.Shipment, error) {

	return s.repository.GetByID(ctx, id)
}

func (s *ShipmentService) GetShipments(
	ctx context.Context,
) ([]*ent.Shipment, error) {

	return s.repository.GetAll(ctx)
}

func (s *ShipmentService) UpdateShipment(
	ctx context.Context,
	id int,
	req model.UpdateShipmentRequest,
) (*ent.Shipment, error) {

	return s.repository.Update(ctx, id, req)
}

func (s *ShipmentService) DeleteShipment(
	ctx context.Context,
	id int,
) error {

	return s.repository.Delete(ctx, id)
}