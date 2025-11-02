package commands

import (
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type AddStoragePlace struct {
	courierID   uuid.UUID
	name        string
	totalVolume int

	isValid bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlace, error) {
	if courierID == uuid.Nil {
		return AddStoragePlace{}, errs.NewValueIsInvalidError("courierID")
	}
	if name == "" {
		return AddStoragePlace{}, errs.NewValueIsInvalidError("name")
	}
	if totalVolume <= 0 {
		return AddStoragePlace{}, errs.NewValueIsInvalidError("totalVolume")
	}

	return AddStoragePlace{
		courierID:   courierID,
		name:        name,
		totalVolume: totalVolume,
		isValid:     true,
	}, nil
}

func (c *AddStoragePlace) IsValid() bool {
	return c.isValid
}

func (c *AddStoragePlace) CourierID() uuid.UUID {
	return c.courierID
}

func (c *AddStoragePlace) Name() string {
	return c.name
}

func (c *AddStoragePlace) TotalVolume() int {
	return c.totalVolume
}
