package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CanDispatchValidOrder(t *testing.T) {
	// Arrange
	couriers := getCouriers()

	orderUuid := uuid.New()
	location, _ := kernel.NewLocation(1, 2)
	newOrder, _ := order.NewOrder(orderUuid, location, 5)

	// Act
	orderService := NewOrderService()
	orderService.Dispatch(newOrder, couriers)
	// Assert
	assert.NotNil(t, newOrder.CourierID(), "Order should be assigned to a courier")

	// Check that the assigned courier is one of the available couriers
	assignedCourierID := *newOrder.CourierID()
	found := false
	for _, courier := range couriers {
		if courier.ID() == assignedCourierID {
			found = true
			break
		}
	}
	assert.True(t, found, "Assigned courier should be one of the available couriers")
}

func getCouriers() []*courier.Courier {
	location1, _ := kernel.NewLocation(1, 1)
	location2, _ := kernel.NewLocation(5, 5)
	location3, _ := kernel.NewLocation(10, 10)

	courier1, _ := courier.NewCourier("Courier 1", 1, location1)
	courier2, _ := courier.NewCourier("Courier 2", 2, location2)
	courier3, _ := courier.NewCourier("Courier 3", 3, location3)

	courier1.AddStoragePlace("storage 1", 10)
	courier2.AddStoragePlace("storage 2", 20)
	courier3.AddStoragePlace("storage 3", 30)

	return []*courier.Courier{courier1, courier2, courier3}
}
