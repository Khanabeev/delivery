package courier

import (
	"delivery/internal/core/domain/kernel"

	"github.com/google/uuid"
)

type Courier struct {
	id           uuid.UUID
	name         string
	speed        int
	location     kernel.Location
	storagePlase []*StoragePlace
}
