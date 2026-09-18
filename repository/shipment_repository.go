package repository

import (
	"context"

	"shipment-service/ent"
	// "shipment-service/ent/shipment"
	"shipment-service/model"
)

type ShipmentRepository struct {
	client *ent.Client
}

func NewShipmentRepository(client *ent.Client) *ShipmentRepository {
	return &ShipmentRepository{
		client: client,
	}
}

func (r *ShipmentRepository) Create(
	ctx context.Context,
	req model.CreateShipmentRequest,
) (*ent.Shipment, error) {

	return r.client.Shipment.
		Create().
		SetTrackingNumber(req.TrackingNumber).
		SetSenderName(req.SenderName).
		SetReceiverName(req.ReceiverName).
		SetOrigin(req.Origin).
		SetDestination(req.Destination).
		SetWeight(req.Weight).
		Save(ctx)
}

func (r *ShipmentRepository) GetByID(
	ctx context.Context,
	id int,
) (*ent.Shipment, error) {

	return r.client.Shipment.
		Get(ctx, id)
}

func (r *ShipmentRepository) GetAll(
	ctx context.Context,
) ([]*ent.Shipment, error) {

	return r.client.Shipment.
		Query().
		All(ctx)
}

func (r *ShipmentRepository) Update(
	ctx context.Context,
	id int,
	req model.UpdateShipmentRequest,
) (*ent.Shipment, error) {

	return r.client.Shipment.
		UpdateOneID(id).
		SetSenderName(req.SenderName).
		SetReceiverName(req.ReceiverName).
		SetOrigin(req.Origin).
		SetDestination(req.Destination).
		SetStatus(req.Status).
		SetWeight(req.Weight).
		Save(ctx)
}

func (r *ShipmentRepository) Delete(
	ctx context.Context,
	id int,
) error {

	return r.client.Shipment.
		DeleteOneID(id).
		Exec(ctx)
}