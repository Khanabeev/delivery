package courierrepo

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CourierRepository = &Repository{}

type Repository struct {
	tracker Tracker
}

func NewRepository(tracker Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &Repository{
		tracker: tracker,
	}, nil
}

// Add implements ports.CourierRepository.
func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	r.tracker.Track(aggregate)

	// DomainToDTO(aggregate) - вызывается функция-маппер,
	// которая преобразует доменную модель в структуру
	// для работы с базой данных.
	dto := DomainToDTO(aggregate)

	// Открыта ли транзакция?
	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	// Вносим изменения
	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	// Если не было внешней в транзакции, то коммитим изменения
	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil

}

// Get implements ports.CourierRepository.
func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	dto := CourierDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, ID)

	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError(ID.String(), uuid.Nil)
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

// GetAll implements ports.CourierRepository.
func (r *Repository) GetAll(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		// Preload all associations
		Preload(clause.Associations).
		Where(`EXISTS (
			SELECT 1 FROM storage_places sp
			WHERE sp.courier_id = couriers.id AND sp.order_id IS NOT NULL
		)`).Find(&dtos)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewObjectNotFoundError("Free courier", nil)
		}
		return nil, result.Error
	}

	aggregates := make([]*courier.Courier, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}
	return aggregates, nil
}

// GetAllFree implements ports.CourierRepository.
func (r *Repository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		// Preload all associations
		Preload(clause.Associations).
		Where(`NOT EXISTS (
			SELECT 1 FROM storage_places sp
			WHERE sp.courier_id = couriers.id AND sp.order_id IS NOT NULL
		)`).Find(&dtos)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewObjectNotFoundError("Free courier", nil)
		}
		return nil, result.Error
	}

	aggregates := make([]*courier.Courier, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}
	return aggregates, nil
}

// Update implements ports.CourierRepository.
func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	r.tracker.Track(aggregate)

	dto := DomainToDTO(aggregate)

	// Открыта ли транзакция?
	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	// Вносим изменения
	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}

	// Если не было внешней в транзакции, то коммитим изменения
	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}
