package commands

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

var _ CreateOrderHandler = &createOrderHandler{}

type createOrderHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateOrderHandler(uowFactory ports.UnitOfWorkFactory) (CreateOrderHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &createOrderHandler{
		uowFactory: uowFactory,
	}, nil
}

func (ch *createOrderHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add create order command")
	}

	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	location := kernel.NewRandomLocation()
	newOrder, err := order.NewOrder(command.OrderID(), location, command.Volume())
	if err != nil {
		return err
	}

	err = uow.OrderRepository().Add(ctx, newOrder)
	if err != nil {
		return err
	}

	return nil
}
