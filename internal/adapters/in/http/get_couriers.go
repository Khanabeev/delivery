package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecasese/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetCouriers(ctx echo.Context) error {
	query := queries.NewGetAllCouriersQuery()

	queryResponses, err := s.getAllCouriersHandler.Handler(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return ctx.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError, problems.NewBadRequest(err.Error()))
	}

	// Map domain responses to HTTP responses
	httpResponses := make([]servers.Courier, len(queryResponses))
	for i, resp := range queryResponses {
		httpResponses[i] = servers.Courier{
			Id:   resp.ID,
			Name: resp.Name,
			Location: servers.Location{
				X: resp.Location.X(),
				Y: resp.Location.Y(),
			},
		}
	}

	return ctx.JSON(http.StatusOK, httpResponses)
}
