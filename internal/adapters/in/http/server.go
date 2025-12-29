package http

import (
	"delivery/internal/core/application/usecasese/commands"
	"delivery/internal/core/application/usecasese/queries"
	"delivery/internal/pkg/errs"
)

type Server struct {
	createOrderCommandHandler      commands.CreateOrderHandler
	createCourierHandler           commands.AddCourierHandler
	getAllCouriersHandler          queries.GetAllCouriersHandler
	getAllUncompletedOrdersHandler queries.GetAllUncompletedOrdersHandler
}

func NewServer(
	createOrderCommandHandler commands.CreateOrderHandler,
	getAllCouriersHandler queries.GetAllCouriersHandler,
	createCourierHandler commands.AddCourierHandler,
	getAllUncompletedOrdersHandler queries.GetAllUncompletedOrdersHandler,
) (*Server, error) {
	if createOrderCommandHandler == nil {

		return nil, errs.NewValueIsRequiredError("createOrderHandler")
	}
	if createCourierHandler == nil {
		return nil, errs.NewValueIsRequiredError("createCourierHandler")
	}
	if getAllCouriersHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriersHandler")
	}
	if getAllUncompletedOrdersHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllUncompletedOrdersHandler")
	}

	return &Server{
		createOrderCommandHandler:      createOrderCommandHandler,
		getAllCouriersHandler:          getAllCouriersHandler,
		createCourierHandler:           createCourierHandler,
		getAllUncompletedOrdersHandler: getAllUncompletedOrdersHandler,
	}, nil
}
