 package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"fmt"
)

var _ OrderDispatcher = &orderService{}

type OrderDispatcher interface {
	Dispatch(*order.Order, []*courier.Courier) (*courier.Courier, error)
}

type orderService struct{}

func NewOrderService() OrderDispatcher {
	return &orderService{}
}

func (os *orderService) Dispatch(newOrder *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {

	if newOrder.Status() != order.StatusCreated {
		message := fmt.Sprintf("Incorrect order status: %v", newOrder.Status())
		return nil, errs.NewValueIsInvalidError(message)
	}

	var fastestCourier *courier.Courier
	var fastestDeliveryTime float64 = -1

	for _, courier := range couriers {
		canTakeOrder, _ := courier.CanTakeOrder(newOrder)

		if canTakeOrder {
			deliveryTime, err := courier.CalculateTimeToLocation(newOrder.Location())
			if err != nil {
				return nil, err
			}

			if fastestDeliveryTime == -1 || deliveryTime < fastestDeliveryTime {
				fastestDeliveryTime = deliveryTime
				fastestCourier = courier
			}
		}
	}

	if fastestCourier == nil {
		return nil, errors.New("Courier not found")
	}

	newOrder.Assign(fastestCourier.ID())
	fastestCourier.TakeOrder(newOrder)

	return fastestCourier, nil
}
