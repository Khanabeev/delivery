package commands

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AssignOrderToCourierHandler interface {
	Handle(context.Context, AssignOrderToCourierCommand) error
}

var _ AssignOrderToCourierHandler = &assignOrderToCourierHandler{}

type assignOrderToCourierHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewAssignOrderToCourierHandler(uowFactory ports.UnitOfWorkFactory) (AssignOrderToCourierHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unit of work Factory")
	}

	return &assignOrderToCourierHandler{
		uowFactory: uowFactory,
	}, nil
}

func (ch *assignOrderToCourierHandler) Handle(ctx context.Context, command AssignOrderToCourierCommand) error {
	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	freeCouriers, err := uow.CourierRepository().GetAllFree(ctx)
	if err != nil {
		return err
	}

	createdOrder, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		return err
	}

	uow.Begin(ctx)
	orderService := services.NewOrderService()
	assignedCourier, err := orderService.Dispatch(createdOrder, freeCouriers)
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Update(ctx, assignedCourier)
	if err != nil {
		return err
	}
	err = uow.OrderRepository().Update(ctx, createdOrder)
	if err != nil {
		return err
	}
	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
