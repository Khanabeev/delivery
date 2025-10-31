package order

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type Order struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	courierID     *uuid.UUID
	location      kernel.Location
	volume        int
	status        Status
}

// ClearDomainEvents implements ddd.AggregateRoot.
func (o *Order) ClearDomainEvents() {
	panic("unimplemented")
}

// GetDomainEvents implements ddd.AggregateRoot.
func (o *Order) GetDomainEvents() []ddd.DomainEvent {
	panic("unimplemented")
}

// RaiseDomainEvent implements ddd.AggregateRoot.
func (o *Order) RaiseDomainEvent(ddd.DomainEvent) {
	panic("unimplemented")
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
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](orderID),
		location:      location,
		volume:        volume,
		status:        StatusCreated,
	}, nil
}

func RestoreOrder(id uuid.UUID, courierID *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		baseAggregate: ddd.NewBaseAggregate(id),
		courierID:     courierID,
		location:      location,
		volume:        volume,
		status:        status,
	}
}

func (o *Order) ID() uuid.UUID {
	return o.baseAggregate.ID()
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
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
