package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AddStoragePlaceHandler interface {
	Handle(context.Context, AddStoragePlaceCommand) error
}

var _ AddStoragePlaceHandler = &addStoragePlaceHandler{}

type addStoragePlaceHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewAddStoragePlaceHandler(uowFactory ports.UnitOfWorkFactory) (AddStoragePlaceHandler, error) {
	if uowFactory == nil {
		return &addStoragePlaceHandler{}, errs.NewValueIsInvalidError("Unit of work Factory is required")
	}

	return &addStoragePlaceHandler{
		uowFactory: uowFactory,
	}, nil
}

func (ch *addStoragePlaceHandler) Handle(ctx context.Context, command AddStoragePlaceCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("add courier command is required")
	}

	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	courierAggregate, err := uow.CourierRepository().Get(ctx, command.courierID)
	if err != nil {
		return err
	}

	err = courierAggregate.AddStoragePlace(command.Name(), command.TotalVolume())
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Update(ctx, courierAggregate)
	if err != nil {
		return err
	}

	return nil
}
