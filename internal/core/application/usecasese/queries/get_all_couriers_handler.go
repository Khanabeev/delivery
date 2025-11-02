package queries

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAllCouriersHandler interface {
	Handler(GetAllCouriersQuery) ([]GetAllCouriersResponse, error)
}

var _ GetAllCouriersHandler = &getAllCouriersHandler{}

type getAllCouriersHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersHandler(db *gorm.DB) (GetAllCouriersHandler, error) {
	if db == nil {
		return &getAllCouriersHandler{}, errs.NewValueIsRequiredError("db")
	}

	return &getAllCouriersHandler{
		db: db,
	}, nil
}

// courierRow используется для сканирования из базы данных
// GORM автоматически мапит snake_case колонки (id, name, location_x, location_y) на CamelCase поля
type courierRow struct {
	ID        uuid.UUID
	Name      string
	LocationX int
	LocationY int
}

func (qh *getAllCouriersHandler) Handler(_ GetAllCouriersQuery) ([]GetAllCouriersResponse, error) {
	var rows []courierRow

	result := qh.db.
		Raw(`
		SELECT id, name, location_x, location_y
		FROM couriers c
		JOIN storage_places sp ON c.id = sp.courier_id AND sp.order_id IS NOT NULL
		ORDER BY id`).
		Scan(&rows)

	if result.Error != nil {
		return nil, result.Error
	}

	// Преобразуем результаты в список GetAllCouriersResponse
	responses := make([]GetAllCouriersResponse, 0, len(rows))
	for _, row := range rows {
		location, err := kernel.NewLocation(row.LocationX, row.LocationY)
		if err != nil {
			return nil, err
		}

		responses = append(responses, GetAllCouriersResponse{
			ID:       row.ID,
			Name:     row.Name,
			Location: location,
		})
	}

	return responses, nil
}
