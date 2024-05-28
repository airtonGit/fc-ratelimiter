//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/airtongit/fullcycle/cleanarch/internal/entity"
	"github.com/airtongit/fullcycle/cleanarch/internal/event"
	"github.com/airtongit/fullcycle/cleanarch/internal/infra/database"
	"github.com/airtongit/fullcycle/cleanarch/internal/infra/web"
	"github.com/airtongit/fullcycle/cleanarch/internal/usecase"
	"github.com/airtongit/fullcycle/cleanarch/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

//var setListedOrderEvent = wire.NewSet(
//	event.NewOrdersListed,
//	wire.Bind(new(events.EventInterface), new(*event.OrdersListed)),
//)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setListedOrderEvent,
		usecase.NewListedOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
