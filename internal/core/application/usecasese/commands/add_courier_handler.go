package commands

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AddCourierHandler interface {
	Handle(context.Context, AddCourierCommand) error
}

var _ AddCourierHandler = &addCourierHandler{}

type addCourierHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewAddCourierHandler(uowFactory ports.UnitOfWorkFactory) (AddCourierHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unit of work Factory")
	}

	return &addCourierHandler{
		uowFactory: uowFactory,
	}, nil
}

func (ch *addCourierHandler) Handle(ctx context.Context, command AddCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("add courier command is required")
	}

	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	location := kernel.NewRandomLocation()
	newCourier, err := courier.NewCourier(command.Name(), command.Speed(), location)
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Add(ctx, newCourier)
	if err != nil {
		return err
	}

	return nil
}
