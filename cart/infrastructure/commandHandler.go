package infrastructure

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type CommandHandler[Aggregate cqrs.Aggregate] struct {
	store cqrs.EventStore
}

func NewCommandHandler[Aggregate cqrs.Aggregate](store cqrs.EventStore) *CommandHandler[Aggregate] {

	return &CommandHandler[Aggregate]{
		store: store,
	}
}

func (h *CommandHandler[Aggregate]) Handle(ctx context.Context, command cqrs.Command) error {
	var aggregate cqrs.Aggregate

	events, err := h.store.LoadFrom(ctx, command.AggregateID(), 0)

	if err != nil {
		return err
	}

	var version uint64

	for event := range events {
		if err := DispatchEvent(aggregate, event); err != nil {
			return err
		}

		version++
	}
	aggregate.SetAggregateVersion(version)

	if err := DispatchCommand(ctx, aggregate, command); err != nil {
		return err
	}

	//events

	if err := h.store.Save(ctx, aggregate.UncommittedEvents(), version); err != nil {
		return err
	}

	return nil

}
