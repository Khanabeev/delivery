package courier

import (
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_StoragePlaceAreCreatedCorrectlyWhenParamsAreCorrect(t *testing.T) {
	// Arrange

	// Act
	StoragePlace, err := NewStoragePlace("backpack", 20)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, StoragePlace)
	assert.NotEmpty(t, StoragePlace.ID())
	assert.Equal(t, "backpack", StoragePlace.Name())
	assert.Empty(t, StoragePlace.OrderID())
	assert.Greater(t, StoragePlace.TotalVolume(), 0)
}

func Test_StoragePlaceShouldReturnErrorWhenParamsAreIncorrectOnCreate(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		name        string
		totalVolume int
		expected    error
	}{
		"wrong_name": {
			name:        "",
			totalVolume: 10,
			expected:    errs.NewValueIsRequiredError("name"),
		},
		"wrong_volume": {
			name:        "backpack",
			totalVolume: -10,
			expected:    errs.NewValueIsInvalidError("totalVolume"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewStoragePlace(test.name, test.totalVolume)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_StoragePlaceCanStore(t *testing.T) {
	// Arrange
	StoragePlace, _ := NewStoragePlace("backpack", 20)
	// Act
	_, err := StoragePlace.CanStore(20)
	// Assert
	assert.NoError(t, err)
}

func Test_StoragePlaceShouldKeepOrder(t *testing.T) {
	// Arrange
	StoragePlace, _ := NewStoragePlace("backpack", 20)
	OrderID := uuid.New()
	// Act
	err := StoragePlace.Store(OrderID, 20)
	// Assert
	assert.NoError(t, err)
}

func Test_StoragePlaceShouldClearOrder(t *testing.T) {
	// Arrange
	StoragePlace, _ := NewStoragePlace("backpack", 20)
	OrderID := uuid.New()
	StoragePlace.Store(OrderID, 20)
	// Act
	err := StoragePlace.Clear(OrderID)
	// Assert
	assert.NoError(t, err)
	assert.Empty(t, StoragePlace.OrderID())
}
