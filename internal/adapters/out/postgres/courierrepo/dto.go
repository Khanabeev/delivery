package courierrepo

import (
	"github.com/google/uuid"
)

type CourierDTO struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Location      LocationDTO `gorm:"embedded;embeddedPrefix:location_"`
	Name          string
	Speed         int
	StoragePlases []StoragePlaceDTO `gorm:"foreignKey:CourierID"`
}

type StoragePlaceDTO struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CourierID   uuid.UUID `gorm:"type:uuid;index"`
	Name        string
	TotalVolume int
	OrderID     *uuid.UUID `gorm:"type:uuid;index"`
}

type LocationDTO struct {
	X int
	Y int
}

func (CourierDTO) TableName() string {
	return "couriers"
}
