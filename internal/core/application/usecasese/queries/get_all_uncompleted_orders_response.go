package queries

import (
	"delivery/internal/core/domain/model/kernel"

	"github.com/google/uuid"
)

type GetAllUncompletedOrdersResponse struct {
	ID       uuid.UUID
	Location kernel.Location
}
