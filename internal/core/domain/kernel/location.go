package kernel

import (
	"delivery/internal/pkg/errs"
	"math"
	"math/rand"
)

type Location struct {
	x     int
	y     int
	isSet bool
}

func NewLocation(x, y int) (Location, error) {
	if x < 1 || x > 10 {
		return Location{}, errs.NewValueIsInvalidError("x")
	}
	if y < 1 || y > 10 {
		return Location{}, errs.NewValueIsInvalidError("y")
	}

	return Location{
		x:     x,
		y:     y,
		isSet: true,
	}, nil
}

func NewRandomLocation() Location {
	min := 1
	max := 10
	x := rand.Intn(max-min) + min
	y := rand.Intn(max-min) + min

	return Location{
		x:     x,
		y:     y,
		isSet: true,
	}
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equal(other Location) bool {
	return l.x == other.x && l.y == other.y
}

func (l Location) IsEmpty() bool {
	return l.isSet
}

func (l Location) DistanceTo(target Location) (int, error) {
	dx := l.x - target.x
	dy := l.y - target.y
	res := int(math.Abs(float64(dx + dy)))

	return res, nil
}
