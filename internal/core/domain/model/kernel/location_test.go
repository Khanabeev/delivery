package kernel

import (
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LocationBeCorrectedWhenParamsAreCorrectOnCreate(t *testing.T) {
	// Arrange

	// Act
	Location, err := NewLocation(1, 2)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, Location)
	assert.Equal(t, 1, Location.X())
	assert.Equal(t, 2, Location.Y())
}

func Test_LocationShouldReturnErrorWhenParametersOutOfRange(t *testing.T) {
	// Arrange

	// Act
	_, err := NewLocation(100, -1)

	assert.Error(t, err)
}

func Test_LocationShouldCalculateDistance(t *testing.T) {
	tests := map[string]struct {
		x1       int
		y1       int
		x2       int
		y2       int
		expected int
	}{
		"check_1": {
			x1:       2,
			y1:       6,
			x2:       4,
			y2:       9,
			expected: 5,
		},
		"check_2": {
			x1:       1,
			y1:       1,
			x2:       1,
			y2:       1,
			expected: 0,
		},

		"check_3": {
			x1:       1,
			y1:       1,
			x2:       10,
			y2:       10,
			expected: 18,
		},

		"check_4": {
			x1:       5,
			y1:       5,
			x2:       6,
			y2:       6,
			expected: 2,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			location1, _ := NewLocation(test.x1, test.y1)
			location2, _ := NewLocation(test.x2, test.y2)
			distance, _ := location1.DistanceTo(location2)

			// Assert
			assert.Equal(t, test.expected, distance)

			if distance != test.expected {
				t.Errorf("expected %v, got %v", test.expected, distance)
			}
		})
	}
}

func Test_LocationReturnErrorWhenParamsAreIncorrectOnCreate(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		x        int
		y        int
		expected error
	}{
		"wrong_x_less_then_1": {
			x:        0,
			y:        1,
			expected: errs.NewValueIsOutOfRangeError("x", 0, minX, maxX),
		},
		"wrong_x_bigger_then_10": {
			x:        11,
			y:        1,
			expected: errs.NewValueIsOutOfRangeError("x", 11, minX, maxX),
		},
		"wrong_y_less_then_1": {
			x:        1,
			y:        0,
			expected: errs.NewValueIsOutOfRangeError("y", 0, minY, maxY),
		},
		"wrong_y_bigger_then_10": {
			x:        1,
			y:        11,
			expected: errs.NewValueIsOutOfRangeError("y", 11, minY, maxY),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewLocation(test.x, test.y)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_LocationShouldCreateRandom(t *testing.T) {
	location := NewRandomLocation()
	assert.GreaterOrEqual(t, location.X(), 1)
	assert.GreaterOrEqual(t, location.Y(), 1)
	assert.LessOrEqual(t, location.X(), 10)
	assert.LessOrEqual(t, location.Y(), 10)
	assert.True(t, location.isValid)
}
