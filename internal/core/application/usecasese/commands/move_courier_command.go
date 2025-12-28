package commands

type MoveCourierCommand struct {
	isValid bool
}

func NewMoveCouriersCommand() (MoveCourierCommand, error) {
	return MoveCourierCommand{
		isValid: true,
	}, nil
}
