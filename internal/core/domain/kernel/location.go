package kernel

import (
	"delivery/internal/pkg/errs"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	minX int = 1
	minY int = 1
	maxX int = 10
	maxY int = 10
)

type Location struct {
	x int
	y int

	valid bool
}

func NewLocation(x, y int) (Location, error) {
	if x < minX || x > maxX {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, minX, maxX)
	}
	if y < minY || y > maxY {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, minY, maxY)
	}

	return Location{
		x: x,
		y: y,

		valid: true,
	}, nil
}

func NewRandomLocation() Location {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	x := r.Intn(maxX-minX) + minX
	y := r.Intn(maxY-minY) + minY

	location, err := NewLocation(x, y)

	if err != nil {
		panic(fmt.Sprintf("invalid random location: x=%d, y=%d, err=%v", x, y, err))
	}

	return location
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

func (l Location) IsValid() bool {
	return l.valid
}

func (l Location) DistanceTo(target Location) (int, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsRequiredError("location")
	}

	dx := l.x - target.x
	dy := l.y - target.y
	res := int(math.Abs(float64(dx + dy)))

	return res, nil
}

func MinLocation() Location {
	location, err := NewLocation(minX, minY)
	if err != nil {
		panic("invalid min location configuration")
	}

	return location
}

func MaxLocation() Location {
	location, err := NewLocation(maxX, maxY)
	if err != nil {
		panic("invalid min location configuration")
	}

	return location
}
