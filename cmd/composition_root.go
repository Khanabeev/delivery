package cmd

import (
	"delivery/internal/adapters/out/postgres"
	commands "delivery/internal/core/application/usecasese/commands"
	queries "delivery/internal/core/application/usecasese/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"fmt"

	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs Config
	db      *gorm.DB

	closers []Closer
}

func NewCompositionRoot(configs Config) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
	}
}

func (f *CompositionRoot) NewOrderService() services.OrderDispatcher {
	orderDispatcher := services.NewOrderService()
	return orderDispatcher
}

func (f *CompositionRoot) NewDB() (*gorm.DB, error) {
	if f.db != nil {
		return f.db, nil
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		f.configs.DbHost,
		f.configs.DbUser,
		f.configs.DbPassword,
		f.configs.DbName,
		f.configs.DbPort,
		f.configs.DbSslMode,
	)

	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	f.db = db
	return db, nil
}

func (f *CompositionRoot) NewUnitOfWorkFactory() (ports.UnitOfWorkFactory, error) {
	db, err := f.NewDB()
	if err != nil {
		return nil, err
	}

	return postgres.NewUnitOfWorkFactory(db)
}

// Query Handlers

func (f *CompositionRoot) NewGetAllCouriersHandler() (queries.GetAllCouriersHandler, error) {
	db, err := f.NewDB()
	if err != nil {
		return nil, err
	}

	return queries.NewGetAllCouriersHandler(db)
}

func (f *CompositionRoot) NewGetAllUncompletedOrdersHandler() (queries.GetAllUncompletedOrdersHandler, error) {
	db, err := f.NewDB()
	if err != nil {
		return nil, err
	}

	return queries.NewGetAllUncompletedOrdersHandler(db)
}

// Command Handlers

func (f *CompositionRoot) NewAddCourierHandler() (commands.AddCourierHandler, error) {
	uowFactory, err := f.NewUnitOfWorkFactory()
	if err != nil {
		return nil, err
	}

	return commands.NewAddCourierHandler(uowFactory)
}

func (f *CompositionRoot) NewAddStoragePlaceHandler() (commands.AddStoragePlaceHandler, error) {
	uowFactory, err := f.NewUnitOfWorkFactory()
	if err != nil {
		return nil, err
	}

	return commands.NewAddStoragePlaceHandler(uowFactory)
}

func (f *CompositionRoot) NewAssignOrderToCourierHandler() (commands.AssignOrderToCourierHandler, error) {
	uowFactory, err := f.NewUnitOfWorkFactory()
	if err != nil {
		return nil, err
	}

	return commands.NewAssignOrderToCourierHandler(uowFactory)
}

func (f *CompositionRoot) NewCreateOrderHandler() (commands.CreateOrderHandler, error) {
	uowFactory, err := f.NewUnitOfWorkFactory()
	if err != nil {
		return nil, err
	}

	return commands.NewCreateOrderHandler(uowFactory)
}

func (f *CompositionRoot) NewMoveCourierHandler() (commands.MoveCourierHandler, error) {
	uowFactory, err := f.NewUnitOfWorkFactory()
	if err != nil {
		return nil, err
	}

	return commands.NewMoveCourierHandler(uowFactory)
}
