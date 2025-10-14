package courier

import (
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

var (
	ErrQuantityIsZeroOrLess = errors.New("Value should not be 0 or less")
)

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if totalVolume <= 0 {
		return nil, errs.NewValueIsInvalidError("totalVolume")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	if other == nil {
		return false
	}
	return s.id == other.id
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume <= 0 {
		return false, ErrQuantityIsZeroOrLess
	}

	if s.totalVolume < volume {
		return false, nil
	}

	if s.isOccupied() {
		return false, nil
	}

	return true, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsRequiredError("orderID")
	}

	if volume <= 0 {
		return errs.NewValueIsRequiredError("volume")
	}

	canStore, err := s.CanStore(volume)
	if err != nil {
		return err
	}

	if !canStore {
		return errors.New("Can't sotre order in this storage place")
	}

	s.orderID = &orderID
	return nil
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsRequiredError("orderID")
	}
	if s.orderID == nil || *s.orderID != orderID {
		return errors.New("Storage is already empty")
	}
	s.orderID = nil
	return nil
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderID != nil
}
