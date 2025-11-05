package commands

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/mocks/core/portsmocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AddStoragePlaceHandlerShouldSuccessWithCorrectProps(t *testing.T) {
	// Arrange
	ctx := context.Background()

	location := kernel.NewRandomLocation()
	courierAggregate, err := courier.NewCourier("Courier1", 2, location)
	assert.NoError(t, err)

	var capturedObj *courier.Courier
	courierRepositoryMock := &portsmocks.CourierRepositoryMock{}
	courierRepositoryMock.
		On("Get", ctx, courierAggregate.ID()).
		Return(courierAggregate, nil).
		Once()

	courierRepositoryMock.
		On("Update", ctx, courierAggregate).
		Run(func(args mock.Arguments) {
			capturedObj = args.Get(1).(*courier.Courier)
		}).
		Return(nil, nil).
		Once()

	unitOfWorkMock := &portsmocks.UnitOfWorkMock{}
	unitOfWorkMock.
		On("CourierRepository").
		Return(courierRepositoryMock)
	unitOfWorkMock.
		On("RollbackUnlessCommitted", ctx).
		Return()

	unitOfWorkFactoryMock := &portsmocks.UnitOfWorkFactoryMock{}
	unitOfWorkFactoryMock.
		On("New", ctx).
		Return(unitOfWorkMock, nil).
		Once()

	// Act
	addStoragePlaceHandler, err := NewAddStoragePlaceHandler(unitOfWorkFactoryMock)
	assert.NoError(t, err)
	addStoragePlaceCommand, err := NewAddStoragePlaceCommand(courierAggregate.ID(), "StoragePlace1", 10)
	assert.NoError(t, err)
	err = addStoragePlaceHandler.Handle(ctx, addStoragePlaceCommand)

	// Assert
	assert.NoError(t, err)
	storagePlaces := capturedObj.StoragePlaces()
	assert.Len(t, storagePlaces, 2, "Expected 2 storage places (backpack + StoragePlace1)")

	found := false
	for _, sp := range storagePlaces {
		if sp.Name() == "StoragePlace1" {
			found = true
			assert.Equal(t, 10, sp.TotalVolume())
			break
		}
	}
	assert.True(t, found, "StoragePlace1 should be in the storage places list")

}
