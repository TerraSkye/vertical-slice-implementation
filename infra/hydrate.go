package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type HydrateHandler interface {
	NewEvent() cqrs.Event
	Apply(ctx context.Context, event cqrs.Event)
}

type genericHydrateHandler[T cqrs.Event] struct {
	handleFunc func(ctx context.Context, event T)
}

// NewEventHandler creates a new EventHandler implementation based on provided function
// and event type inferred from function argument.
func NewHydrateHandler[T cqrs.Event](
	handleFunc func(ctx context.Context, event T),
) HydrateHandler {
	return &genericHydrateHandler[T]{
		handleFunc: handleFunc,
	}
}

func (c genericHydrateHandler[T]) NewEvent() cqrs.Event {
	tVar := new(T)
	return *tVar
}

func (c genericHydrateHandler[T]) Apply(ctx context.Context, e cqrs.Event) {
	event := e.(T)
	c.handleFunc(ctx, event)
}

func Hydrate(handlers ...HydrateHandler) func(ctx context.Context, ev cqrs.Event) {
	eventHandlers := make(map[string]HydrateHandler)

	for _, handler := range handlers {
		eventHandlers[cqrs.TypeName(handler.NewEvent())] = handler
	}

	return func(ctx context.Context, ev cqrs.Event) {
		eventName := cqrs.TypeName(ev)
		if handler, ok := eventHandlers[eventName]; ok {
			handler.Apply(ctx, ev)
		}
	}
}
