package courierrepo

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO
	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed()

	// Convert domain StoragePlaces to DTOs
	storagePlaces := aggregate.StoragePlaces()
	storagePlaceDTOs := make([]StoragePlaceDTO, len(storagePlaces))
	for i, sp := range storagePlaces {
		storagePlaceDTOs[i] = StoragePlaceDTO{
			ID:          sp.ID(),
			CourierID:   aggregate.ID(),
			Name:        sp.Name(),
			TotalVolume: sp.TotalVolume(),
			OrderID:     sp.OrderID(),
		}
	}
	courierDTO.StoragePlases = storagePlaceDTOs

	courierDTO.Location = LocationDTO{
		X: aggregate.Location().X(),
		Y: aggregate.Location().Y(),
	}

	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)

	// Convert DTO StoragePlaces to domain StoragePlaces
	storagePlaces := make([]*courier.StoragePlace, len(dto.StoragePlases))
	for i, spDTO := range dto.StoragePlases {
		storagePlaces[i] = courier.RestoreStoragePlace(
			spDTO.ID,
			spDTO.Name,
			spDTO.TotalVolume,
			spDTO.OrderID,
		)
	}

	aggregate := courier.RestoreCourier(dto.ID, dto.Name, location, dto.Speed, storagePlaces)
	return aggregate
}
