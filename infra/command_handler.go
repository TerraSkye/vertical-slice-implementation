package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type CommandHandler struct {
	store cqrs.EventStore
}

// NewCommandHandler creates a new CommandHandler for an aggregate type.
func NewCommandHandler(store cqrs.EventStore) *CommandHandler {
	h := &CommandHandler{
		store: store,
	}

	return h
}

func (h *CommandHandler) Handle(ctx context.Context, command cqrs.Command) error {

	aggregateType, err := cart.AggregateForCommand(command)

	if err != nil {
		return err
	}

	aggregate := aggregateType.New(command.AggregateID())

	events, err := h.store.LoadFrom(ctx, command.AggregateID(), 0)

	if err != nil {
		return err
	}

	var version uint64

	for event := range events {
		if err := cart.DispatchEvent(aggregate, event.Event); err != nil {
			return err
		}

		version++
	}
	aggregate.SetAggregateVersion(version)

	if err := cart.DispatchCommand(ctx, aggregate, command); err != nil {
		return err
	}

	//events
	uncomittedEvents := aggregate.UncommittedEvents()

	if err := h.store.Save(ctx, uncomittedEvents, version); err != nil {
		return err
	}

	return nil

}
