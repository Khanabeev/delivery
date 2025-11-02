package queries

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAllUncompletedOrdersHandler interface {
	Handle(context.Context, GetAllUncompletedOrdersQuery) ([]GetAllUncompletedOrdersResponse, error)
}

var _ GetAllUncompletedOrdersHandler = &getAllUncompletedOrdersHandler{}

type getAllUncompletedOrdersHandler struct {
	db *gorm.DB
}

func NewGetAllUncompletedOrdersHandler(db *gorm.DB) (GetAllUncompletedOrdersHandler, error) {
	if db == nil {
		return &getAllUncompletedOrdersHandler{}, errs.NewValueIsRequiredError("db")
	}

	return &getAllUncompletedOrdersHandler{
		db: db,
	}, nil
}

type orderRow struct {
	ID        uuid.UUID
	LocationX int
	LocationY int
}

func (qh *getAllUncompletedOrdersHandler) Handle(ctx context.Context, _ GetAllUncompletedOrdersQuery) ([]GetAllUncompletedOrdersResponse, error) {
	var rows []orderRow

	result := qh.db.WithContext(ctx).
		Raw(`
			SELECT id, location_x, location_y
			FROM orders
			WHERE status IN (?, ?)
			ORDER BY id
		`, order.StatusCreated, order.StatusAssigned).
		Scan(&rows)

	if result.Error != nil {
		return nil, result.Error
	}

	// Преобразуем результаты в список GetAllUncompletedOrdersResponse
	responses := make([]GetAllUncompletedOrdersResponse, 0, len(rows))
	for _, row := range rows {
		location, err := kernel.NewLocation(row.LocationX, row.LocationY)
		if err != nil {
			return nil, err
		}

		responses = append(responses, GetAllUncompletedOrdersResponse{
			ID:       row.ID,
			Location: location,
		})
	}

	return responses, nil
}
