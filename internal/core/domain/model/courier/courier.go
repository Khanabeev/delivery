package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"math"

	"github.com/google/uuid"
)

type Courier struct {
	id            uuid.UUID
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

	backpack, err := NewStoragePlace("backpack", 10)
	if err != nil {
		return nil, err
	}
	var storagePlases []*StoragePlace

	storagePlases = append(storagePlases, backpack)

	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlases: storagePlases,
	}, nil
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}

	return c.id == other.id
}

func (c *Courier) ID() uuid.UUID {
	return c.id
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

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlases
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
		return false, errs.NewValueIsInvalidError("order")
	}

	for _, s := range c.storagePlases {
		if s.TotalVolume() >= order.Volumne() {
			return true, nil
		}
	}

	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	canTake, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if canTake {
		order.Assign(c.id)
		for _, s := range c.storagePlases {
			if s.TotalVolume() <= order.Volumne() {
				s.orderID = order.ID()
			}
		}
	}

	return nil
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsInvalidError("order")
	}

	isOrderFound := false
	for _, s := range c.StoragePlaces() {
		if s.orderID == order.ID() {
			isOrderFound = true
			s.orderID = nil
		}
	}

	if !isOrderFound {
		return errs.NewObjectNotFoundError("order", order.ID())
	}

	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsRequiredError("target")
	}

	distance, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}

	timeToLocation := math.Ceil((float64)(distance / c.speed))
	return timeToLocation, nil
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
