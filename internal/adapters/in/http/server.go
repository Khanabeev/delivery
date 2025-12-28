package http

import (
	"delivery/internal/core/application/usecasese/commands"
	"delivery/internal/pkg/errs"
)

type Server struct {
	createOrderCommandHandler commands.CreateOrderHandler
}

func NewServer(
	createOrderCommandHandler commands.CreateOrderHandler,
) (*Server, error) {
	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderHandler")
	}

	return &Server{
		createOrderCommandHandler: createOrderCommandHandler,
	}, nil
}
