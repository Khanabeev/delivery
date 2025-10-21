package order

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type Order struct {
	baseEntity *ddd.BaseEntity[uuid.UUID]
	courierID  *uuid.UUID
	location   kernel.Location
	volume     int
	status     Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if !location.IsValid() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		baseEntity: ddd.NewBaseEntity(uuid.New()),
		location:   location,
		volume:     volume,
		status:     StatusCreated,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.baseEntity.ID()
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volumne() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if o.status == StatusAssigned {
		return errs.NewValueIsInvalidError("Courier is already assigned")
	}

	o.courierID = &courierID
	o.status = StatusAssigned

	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return errs.NewValueIsInvalidError("Status should be assigned")
	}

	o.status = StatusCompleted

	return nil
}
