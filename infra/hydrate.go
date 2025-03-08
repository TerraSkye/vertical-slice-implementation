package infra

import (
	"fmt"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"strings"
)

type HydrateHandler interface {
	NewEvent() cqrs.Event
	Apply(event cqrs.Event)
}

type genericHydrateHandler[T cqrs.Event] struct {
	handleFunc func(event T)
}

// NewEventHandler creates a new EventHandler implementation based on provided function
// and event type inferred from function argument.
func NewHydrateHandler[T cqrs.Event](
	handleFunc func(event T),
) HydrateHandler {
	return &genericHydrateHandler[T]{
		handleFunc: handleFunc,
	}
}

func (c genericHydrateHandler[T]) NewEvent() cqrs.Event {
	tVar := new(T)
	return *tVar
}

func (c genericHydrateHandler[T]) Apply(e cqrs.Event) {
	event := e.(T)
	c.handleFunc(event)
}

func Hydrate(handlers ...HydrateHandler) func(ev cqrs.Event) {
	eventHandlers := make(map[string]HydrateHandler)

	// StructName name returns struct name in format [type name].
	// It ignores if the value is a pointer or not.
	structName := func(v interface{}) string {
		segments := strings.Split(fmt.Sprintf("%T", v), ".")

		return segments[len(segments)-1]
	}

	for _, handler := range handlers {
		eventHandlers[structName(handler.NewEvent())] = handler
	}

	return func(ev cqrs.Event) {
		eventName := structName(ev)
		if handler, ok := eventHandlers[eventName]; ok {
			handler.Apply(ev)
		}
	}
}
