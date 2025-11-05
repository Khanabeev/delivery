package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type MoveCourierHandler interface {
	Handler(context.Context, MoveCourierCommand) error
}

var _ MoveCourierHandler = &moveCourierHandler{}

type moveCourierHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewMoveCourierHandler(uowFactory ports.UnitOfWorkFactory) (MoveCourierHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsInvalidError("unit of work factory is required")
	}

	return &moveCourierHandler{
		uowFactory: uowFactory,
	}, nil
}

func (ch *moveCourierHandler) Handler(ctx context.Context, command MoveCourierCommand) error {
	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	courierAggregates, err := uow.CourierRepository().GetAll(ctx)
	if err != nil {
		return err
	}

	for _, courierAgg := range courierAggregates {
		for _, sp := range courierAgg.StoragePlaces() {
			if sp.OrderID() != nil {
				orderAggregate, err := uow.OrderRepository().Get(ctx, *sp.OrderID())
				if err != nil {
					return err
				}
				if orderAggregate != nil {
					// Перемещаем курьера на 1 шаг в сторону заказа со скоростью его транспорта
					err = courierAgg.Move(orderAggregate.Location())
					if err != nil {
						return err
					}

					// Проверяем, совпадают ли координаты курьера и заказа
					if courierAgg.Location().Equal(orderAggregate.Location()) {
						// Завершаем заказ (переводим в Completed)
						err = orderAggregate.Complete()
						if err != nil {
							return err
						}

						// Курьер освобождает место хранения
						err = courierAgg.CompleteOrder(orderAggregate)
						if err != nil {
							return err
						}

						// Сохраняем изменения заказа
						err = uow.OrderRepository().Update(ctx, orderAggregate)
						if err != nil {
							return err
						}
					}

					// Сохраняем изменения курьера
					err = uow.CourierRepository().Update(ctx, courierAgg)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Коммитим транзакцию
	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
