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

func Test_AddCourierHandlerShouldSuccessWhenParamsAreCorrect(t *testing.T) {
	//Arrange
	ctx := context.Background()
	location := kernel.NewRandomLocation()

	courierAggregate, err := courier.NewCourier("Courier1", 2, location)
	assert.NoError(t, err)

	courierRepositoryMock := &portsmocks.CourierRepositoryMock{}
	courierRepositoryMock.
		On("Add", ctx, mock.MatchedBy(func(orderArg *courier.Courier) bool {
			return courierAggregate.Name() == "Courier1" && courierAggregate.Speed() == 2
		})).
		Return(nil).
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

	//Act
	addCourierHandler, err := NewAddCourierHandler(unitOfWorkFactoryMock)
	assert.NoError(t, err)
	addCourierCommand, err := NewAddNewCourierCommand("Courier1", 2)
	assert.NoError(t, err)
	err = addCourierHandler.Handle(ctx, addCourierCommand)

	//Assert
	assert.NoError(t, err)
}
