package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/io-da/query"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/additem"
	"github.com/terraskye/vertical-slice-implementation/cart/archiveitem"
	"github.com/terraskye/vertical-slice-implementation/cart/cartitems"
	"github.com/terraskye/vertical-slice-implementation/cart/cartwithproducts"
	_ "github.com/terraskye/vertical-slice-implementation/cart/handlers"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"github.com/terraskye/vertical-slice-implementation/infra"
	"net/http"
)

func main() {

	setupOTelSDK(context.Background())

	var store cqrs.EventStore
	var router = mux.NewRouter()

	var eventBus infra.EventBus

	{
		// the bus is in memory
		eventBus = infra.NewEventBus()
	}

	{
		store = infra.NewMemoryStore(eventBus)
	}

	var queryBus *query.Bus

	var queryProvider infra.QueryProvider
	var queryIteratorProvider infra.QueryIteratorProvider

	{
		// the query bus
		queryBus = query.NewBus()

		queryProvider = infra.NewQueryHandler()
		queryIteratorProvider = infra.NewQueryIteratorHandler()

		queryBus.Handlers(queryProvider)
		queryBus.InitializeIteratorHandlers(queryIteratorProvider)

	}

	var commandBus infra.CommandBus

	{
		commandBus = infra.NewCommandBus(20)
		//a command handler per aggregate type?
		commandBus.AddHandler(infra.NewCommandHandler(store,
			cart.AggregateForCommand,
			cart.DispatchEvent,
			cart.DispatchCommand,
		).Handle)
	}

	{
		// additem
		service := additem.NewService(commandBus)
		additem.MakeHttpHandler(router, service)
	}

	{
		service := archiveitem.NewService(commandBus)
		archiveitem.MakeHttpHandler(router, service)

		automation := archiveitem.NewAutomation(commandBus, queryBus)

		// register on the eventbus
		eventBus.Subscribe(infra.NewEventHandler("archiveitem", automation.OnPriceChanged))
		//TODO register on eventbus

	}

	{
		queryHandler := cartitems.NewQueryHandler(store)
		//queryHandler.HandleQuery()
		_ = queryHandler
		queryBus.Handlers()
		//queryBus.Handlers(queryHandler.)
		//TODO register on the query bus.
		cartitems.MakeHttpHandler(router, queryBus)
	}

	{
		projector := cartwithproducts.NewProjector()

		eventBus.SubscribeToGroup(infra.NewEventGroupProcessor("cartwithproducts",
			infra.NewGroupEventHandler(projector.OnItemAdded),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnCartCreated),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnItemRemoved),
		))

		eventBus.SubscribeToGroup(infra.NewEventGroupProcessor("cartwithproducts3",
			infra.NewGroupEventHandler(projector.OnItemAdded),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnCartCreated),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnItemRemoved),
		))

		eventBus.SubscribeToGroup(infra.NewEventGroupProcessor("cartwithproducts4",
			infra.NewGroupEventHandler(projector.OnItemAdded),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnCartCreated),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnItemRemoved),
		))

		//eventBus.Subscribe()

		//TODO register this onto the BUS
		queryHandler := cartwithproducts.NewQueryHandler()

		queryProvider.RegisterHandler(queryHandler)

	}

	http.Handle("/", router)

	errs := make(chan error, 2)

	go func() {
		httpAddr := ":9090"
		fmt.Println("serving on 0.0.0.0:9090")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()

	fmt.Println(<-errs)
}
