package commands

import "delivery/internal/pkg/errs"

type AddCourierCommand struct {
	name    string
	speed   int
	isValid bool
}

func NewAddNewCourierCommand(name string, speed int) (AddCourierCommand, error) {
	if name == "" {
		return AddCourierCommand{}, errs.NewValueIsInvalidError("name")
	}
	if speed < 0 {
		return AddCourierCommand{}, errs.NewValueIsInvalidError("speed")
	}

	return AddCourierCommand{
		name:    name,
		speed:   speed,
		isValid: true,
	}, nil
}

func (c *AddCourierCommand) IsValid() bool {
	return c.isValid
}

func (c *AddCourierCommand) Name() string {
	return c.name
}

func (c *AddCourierCommand) Speed() int {
	return c.speed
}
