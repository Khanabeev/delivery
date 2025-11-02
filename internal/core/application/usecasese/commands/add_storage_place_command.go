package commands

import (
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type AddStoragePlaceCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int

	isValid bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlaceCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("courierID")
	}
	if name == "" {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("name")
	}
	if totalVolume <= 0 {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("totalVolume")
	}

	return AddStoragePlaceCommand{
		courierID:   courierID,
		name:        name,
		totalVolume: totalVolume,
		isValid:     true,
	}, nil
}

func (c *AddStoragePlaceCommand) IsValid() bool {
	return c.isValid
}

func (c *AddStoragePlaceCommand) CourierID() uuid.UUID {
	return c.courierID
}

func (c *AddStoragePlaceCommand) Name() string {
	return c.name
}

func (c *AddStoragePlaceCommand) TotalVolume() int {
	return c.totalVolume
}
