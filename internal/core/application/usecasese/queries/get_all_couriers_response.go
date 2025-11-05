package queries

import (
	"delivery/internal/core/domain/model/kernel"

	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	ID       uuid.UUID
	Name     string
	Location kernel.Location
}
