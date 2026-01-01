package commands

type AssignOrderToCourierCommand struct {
	isValid bool
}

func NewAssignOrderToCourierCommand() (AssignOrderToCourierCommand, error) {
	return AssignOrderToCourierCommand{
		isValid: true,
	}, nil
}
