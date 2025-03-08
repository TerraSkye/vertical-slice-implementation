package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/io-da/query"
	"github.com/terraskye/vertical-slice-implementation/cart/additem"
	"github.com/terraskye/vertical-slice-implementation/cart/archiveitem"
	"github.com/terraskye/vertical-slice-implementation/cart/cartitems"
	"github.com/terraskye/vertical-slice-implementation/cart/cartwithproducts"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	_ "github.com/terraskye/vertical-slice-implementation/cart/handlers"
	"github.com/terraskye/vertical-slice-implementation/cart/infrastructure"
	_ "github.com/terraskye/vertical-slice-implementation/cart/infrastructure"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"github.com/terraskye/vertical-slice-implementation/infra"
	"net/http"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {

	var store cqrs.EventStore
	var router = mux.NewRouter()

	{
		store = infra.NewMemoryStore()
	}

	var queryBus *query.Bus

	{
		queryBus = query.NewBus()
	}

	var eventBus infra.EventBus

	{

	}

	var commandBus infra.CommandBus

	{

		commandBus.AddHandler(infrastructure.NewCommandHandler[domain.Cart](store).Handle)
		commandBus.AddHandler(infrastructure.NewCommandHandler[domain.Pricing](store).Handle)
		commandBus.AddHandler(infrastructure.NewCommandHandler[domain.Inventory](store).Handle)
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
		//queryBus.Handlers(queryHandler)
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
		httpAddr := ":8080"
		fmt.Println("serving on 0.0.0.0:8080")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()

	fmt.Println(<-errs)
}
