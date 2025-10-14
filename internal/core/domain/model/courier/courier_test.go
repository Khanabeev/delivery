package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CanCrateCourier(t *testing.T) {
	// Arrange
	location := kernel.NewRandomLocation()
	// Act
	courier, err := NewCourier("Jhon", 2, location)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, courier)
	assert.Equal(t, 1, len(courier.storagePlases))

	for _, v := range courier.storagePlases {
		assert.Equal(t, 10, v.TotalVolume())
	}
}

func Test_CanAddNewStorage(t *testing.T) {
	// Arrange
	courier := getNewCourier()
	// Act
	err := courier.AddStoragePlace("trunk", 100)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(courier.storagePlases))
	for _, v := range courier.storagePlases {
		if v.Name() == "trunk" {
			assert.Equal(t, 100, v.TotalVolume())
		}
	}
}

func Test_CourierCanTakeOrder(t *testing.T) {
	// Arrange
	courier := getNewCourier()
	location := kernel.NewRandomLocation()

	orderId := uuid.New()
	newOrder1, err := order.NewOrder(orderId, location, 10)
	assert.NoError(t, err)

	orderId2 := uuid.New()
	newOrder2, err := order.NewOrder(orderId2, location, 100)
	assert.NoError(t, err)

	// Act
	canTake1, err := courier.CanTakeOrder(newOrder1)
	assert.NoError(t, err)

	canTake2, err := courier.CanTakeOrder(newOrder2)
	assert.NoError(t, err)
	// Assert

	assert.NotEmpty(t, newOrder1)
	assert.NotEmpty(t, newOrder2)
	assert.True(t, canTake1)
	assert.False(t, canTake2)
}

func getNewCourier() *Courier {
	location := kernel.NewRandomLocation()
	courier, _ := NewCourier("Jhon", 2, location)
	return courier
}
