package service

import (
	"context"
	"fmt"

	"shipment-service/ent"
	"shipment-service/model"
	"shipment-service/repository"
)

type ShipmentService struct {
	repository       *repository.ShipmentRepository
	searchRepository *repository.ShipmentSearchRepository
}

func NewShipmentService(
	repository *repository.ShipmentRepository,
	searchRepository *repository.ShipmentSearchRepository,
) *ShipmentService {
	return &ShipmentService{
		repository:       repository,
		searchRepository: searchRepository,
	}
}

func (s *ShipmentService) SearchShipments(
	ctx context.Context,
	query string,
) ([]model.ShipmentDocument, error) {
	return s.searchRepository.Search(ctx, query)
}

func (s *ShipmentService) CreateShipment(
	ctx context.Context,
	req model.CreateShipmentRequest,
) (*ent.Shipment, error) {
	shipment, err := s.repository.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.searchRepository.Index(ctx, shipment); err != nil {
		return nil, fmt.Errorf("sync created shipment with search index: %w", err)
	}

	return shipment, nil
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
	shipment, err := s.repository.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	if err := s.searchRepository.Index(ctx, shipment); err != nil {
		return nil, fmt.Errorf("sync updated shipment with search index: %w", err)
	}

	return shipment, nil
}

func (s *ShipmentService) DeleteShipment(
	ctx context.Context,
	id int,
) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.searchRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("sync deleted shipment with search index: %w", err)
	}

	return nil
}
