package courier

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

type Courier struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	name          string
	speed         int
	location      kernel.Location
	storagePlases []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsInvalidError("name")
	}

	if speed <= 0 {
		return nil, errs.NewValueIsInvalidError("speed")
	}

	courier := &Courier{
		baseAggregate: ddd.NewBaseAggregate(uuid.New()),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlases: make([]*StoragePlace, 0),
	}

	err := courier.AddStoragePlace("backpack", 10)
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func RestoreCourier(id uuid.UUID, name string, location kernel.Location, speed int, storagePlaces []*StoragePlace) *Courier {
	return &Courier{
		baseAggregate: ddd.NewBaseAggregate(id),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlases: storagePlaces,
	}
}

func (c *Courier) ID() uuid.UUID {
	return c.baseAggregate.ID()
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []StoragePlace {
	res := make([]StoragePlace, len(c.storagePlases))
	for i, sp := range c.storagePlases {
		res[i] = *sp
	}

	return res
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	newStoragePlace, err := NewStoragePlace(name, volume)

	if err != nil {
		return err
	}

	c.storagePlases = append(c.storagePlases, newStoragePlace)

	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, errs.NewValueIsRequiredError("order")
	}

	for _, s := range c.storagePlases {
		canStore, err := s.CanStore(order.Volume())
		if err != nil {
			return false, err
		}

		if canStore {
			return true, nil
		}
	}

	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {

	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	canTake, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if !canTake {
		return errors.New("No storage place")
	}

	for _, s := range c.storagePlases {
		canStore, err := s.CanStore(order.Volume())
		if err != nil {
			return err
		}

		if canStore {
			err := s.Store(order.ID(), order.Volume())
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("No storage place")
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsInvalidError("order")
	}

	storagePlace, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return errs.NewObjectNotFoundError("order", order.ID())
	}

	err = storagePlace.Clear(order.ID())
	if err != nil {
		return err
	}

	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	for _, s := range c.storagePlases {
		if s.orderID == &orderID {
			return s, nil
		}
	}
	return nil, errors.New("Storage place not found")
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsRequiredError("target")
	}

	distance, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}

	time := float64(distance) / float64(c.speed)
	return time, nil
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsValid() {
		return errs.NewValueIsRequiredError("target")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) IsFree() bool {
	for _, storagePlace := range c.storagePlases {
		if storagePlace.OrderID() != nil {
			return false
		}
	}
	return true
}

// DDD AggregateRoot interface implementation
func (c *Courier) GetDomainEvents() []ddd.DomainEvent {
	return c.baseAggregate.GetDomainEvents()
}

func (c *Courier) ClearDomainEvents() {
	c.baseAggregate.ClearDomainEvents()
}

func (c *Courier) RaiseDomainEvent(event ddd.DomainEvent) {
	c.baseAggregate.RaiseDomainEvent(event)
}
