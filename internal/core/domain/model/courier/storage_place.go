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
		return nil, ErrQuantityIsZeroOrLess
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
		orderID:     nil,
	}, nil
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
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
		return false, errors.New("Total Volume is less that requested volume")
	}

	if s.isOccupied() {
		return false, errors.New("Storage is occupied")
	}

	return true, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	_, err := s.CanStore(volume)
	if err != nil {
		return err
	}

	s.orderID = &orderID
	return nil
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if s.orderID == nil {
		return errors.New("Storage is already empty")
	}
	s.orderID = nil
	return nil
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderID != nil
}
