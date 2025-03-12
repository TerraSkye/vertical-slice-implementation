package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/io-da/query"
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

	var store cqrs.EventStore
	var router = mux.NewRouter()

	var eventBus infra.EventBus

	{
		eventBus = infra.NewEventBus()
	}

	{
		store = infra.NewMemoryStore(eventBus)
	}

	var queryBus *query.Bus

	{
		queryBus = query.NewBus()
	}

	var commandBus infra.CommandBus

	{
		commandBus = infra.NewCommandBus(20)
		commandBus.AddHandler(infra.NewCommandHandler(store).Handle)
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

		_ = queryHandler
		//queryBus.Handlers(queryHandler.)
		//TODO register on the query bus.
		cartitems.MakeHttpHandler(router, queryBus)
	}

	{
		projector := cartwithproducts.NewProjector()

		infra.NewEventGroupProcessor(
			infra.NewGroupEventHandler(projector.OnItemAdded),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnCartCreated),
			infra.NewGroupEventHandler(projector.OnItemArchived),
			infra.NewGroupEventHandler(projector.OnItemRemoved),
		)

		//eventBus.Subscribe()

		//TODO register this onto the BUS
		queryHandler := cartwithproducts.NewQueryHandler()
		_ = queryHandler

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
