package order

import (
	"delivery/internal/core/domain/kernel"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CanCreateOrder(t *testing.T) {
	// Arrange
	orderId := uuid.New()
	location := kernel.NewRandomLocation()
	// Act
	order, err := NewOrder(orderId, location, 15)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
}
