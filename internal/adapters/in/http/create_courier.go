package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecasese/commands"
	"delivery/internal/generated/servers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateCourier(ctx echo.Context) error {
	var courier servers.Courier
	err := ctx.Bind(courier)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	addNewCourierCommand, err := commands.NewAddNewCourierCommand(courier.Name, 1)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createCourierHandler.Handle(ctx.Request().Context(), addNewCourierCommand)
	if err != nil {
		return problems.NewConflict(err.Error(), "/")
	}

	return ctx.JSON(http.StatusOK, nil)
}
