package commands

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/mocks/core/portsmocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateOrderHandlerShouldSuccessWhenParamsAreCorrect(t *testing.T) {
	// Arrange
	ctx := context.Background()
	location := kernel.NewRandomLocation()

	orderID := uuid.New()
	orderAggregate, err := order.NewOrder(orderID, location, 10)
	assert.NoError(t, err)

	orderRepositoryMock := &portsmocks.OrderRepositoryMock{}
	orderRepositoryMock.
		On("Add", ctx, mock.MatchedBy(func(orderArg *order.Order) bool {
			return orderArg.ID() == orderID && orderArg.Volume() == 10
		})).
		Return(nil).
		Once()
	unitOfWorkMock := &portsmocks.UnitOfWorkMock{}
	unitOfWorkMock.
		On("OrderRepository").
		Return(orderRepositoryMock)
	unitOfWorkMock.
		On("RollbackUnlessCommitted", ctx).
		Return()

	unitOfWorkFactoryMock := &portsmocks.UnitOfWorkFactoryMock{}
	unitOfWorkFactoryMock.
		On("New", ctx).
		Return(unitOfWorkMock, nil).
		Once()

	// Act
	createOrderHandler, err := NewCreateOrderHandler(unitOfWorkFactoryMock)
	assert.NoError(t, err)
	createOrderCommand, err := NewCreateOrderCommand(orderAggregate.ID(), "Street", 10)
	assert.NoError(t, err)
	err = createOrderHandler.Handle(ctx, createOrderCommand)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orderID, orderAggregate.ID())
	assert.Equal(t, 10, orderAggregate.Volume())
}
