package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// Command and event registries store registered handlers.
var (
	commandRegistry   = make(map[string]any)                           // Stores command handlers
	eventRegistry     = make(map[string]any)                           // Stores event handlers
	aggregateRegistry = make(map[string]func(id uuid.UUID) any)        // Stores aggregate constructors
	eventDecoder      = make(map[string]func(raw []byte) (any, error)) // Stores event decoders
)

// CommandHandler defines a function type for handling commands.
// It takes an aggregate and returns a function that processes the given command.
type CommandHandler[A any, T any] func(aggregate A) func(ctx context.Context, command T) error

// EventHandler defines a function type for handling events.
// It takes an aggregate and returns a function that processes the given event.
type EventHandler[A any, T cqrs.Event] func(aggregate A) func(event T)

// AggregateHandler defines a function type for instantiating an aggregate.
// It takes an id and returns a function that returns an instantiating aggregate.
type AggregateHandler[A any] func(id uuid.UUID) A

// RegisterCommand registers a command handler for a given aggregate and command type.
// It associates the command type with its handler and ensures the aggregate registry
// can instantiate the appropriate aggregate.
//
// Example usage:
//
//	func init() {
//	    RegisterCommand(func(aggregate *domain.Cart) func(ctx context.Context, command *commands.AddItem) error {
//	        return aggregate.AddItem
//	    })
//	}
func RegisterCommand[A cqrs.Aggregate, T cqrs.Command](handler CommandHandler[A, T]) {
	var cmd T
	cmdType := cqrs.TypeName(cmd)

	// Store a type-erased function that preserves the correct signature
	commandRegistry[cmdType] = func(aggregate cqrs.Aggregate) func(context.Context, cqrs.Command) error {
		return func(ctx context.Context, command cqrs.Command) error {
			return handler(aggregate.(A))(ctx, command.(T)) // Explicit type assertion
		}
	}
}

// RegisterEvent registers an event handler for a given aggregate and event type.
// It associates the event type with its handler and sets up a decoder for the event.
//
// Example usage:
//
//	func init() {
//	    RegisterEvent(func(aggregate *domain.Cart) func(event *events.CartCleared) {
//	        return aggregate.OnCartCleared
//	    })
//	}
func RegisterEvent[A cqrs.Aggregate, T cqrs.Event](handler EventHandler[A, T]) {
	var evt T
	evtType := cqrs.TypeName(evt)
	eventRegistry[evtType] = handler

	// Store an event decoder
	eventDecoder[evtType] = func(raw []byte) (any, error) {
		var evt T
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return evt, nil
	}
}

// RegisterCommand registers a command handler for a given aggregate and command type.
// It associates the command type with its handler and ensures the aggregate registry
// can instantiate the appropriate aggregate.
//
// Example usage:
//
//func init() {
//	RegisterAggregate(func(id uuid.UUID) *domain.Cart {
//		return &domain.Cart{
//			AggregateBase: infra.NewAggregateBase(id),
//		}
//	})
//}

func RegisterAggregate[A cqrs.Aggregate](handler AggregateHandler[A]) {
	var evt A
	aggType := cqrs.TypeName(evt)
	aggregateRegistry[aggType] = func(id uuid.UUID) any {
		return handler(id)
	}
}

// DispatchCommand executes the registered handler for a given command.
// It retrieves the appropriate handler based on the command type and invokes it.
// Returns an error if no handler is registered or if there's a type mismatch.
func DispatchCommand[A cqrs.Aggregate, T cqrs.Command](ctx context.Context, aggregate A, command T) error {
	cmdType := cqrs.TypeName(command)

	handlerRaw, exists := commandRegistry[cmdType]
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", cmdType)
	}

	// Ensure type safety
	handlerWrapper, ok := handlerRaw.(func(cqrs.Aggregate) func(context.Context, cqrs.Command) error)
	if !ok {
		return fmt.Errorf("invalid command handler type for: %s", cmdType)
	}

	// Execute the command handler
	return handlerWrapper(aggregate)(ctx, command)
}

// DispatchEvent executes the registered handler for a given event.
// It retrieves the appropriate handler based on the event type and invokes it.
// Returns an error if no handler is registered or if there's a type mismatch.
func DispatchEvent[A cqrs.Aggregate, T cqrs.Event](aggregate A, event T) error {
	evtType := cqrs.TypeName(event)
	handlerRaw, exists := eventRegistry[evtType]
	if !exists {
		return errors.New("no handler registered for event: " + evtType)
	}

	// Ensure type safety
	handler, ok := handlerRaw.(EventHandler[A, T])
	if !ok {
		return errors.New("invalid event handler type for: " + evtType)
	}

	// Execute the event handler
	handler(aggregate)(event)
	return nil
}

// AggregateForCommand retrieves the appropriate aggregate for a given command.
// It uses the command type to find and instantiate the corresponding aggregate.
// Returns an error if no matching aggregate is found.
func AggregateForCommand(cmd cqrs.Command) (cqrs.Aggregate, error) {
	cmdType := cqrs.TypeName(cmd)

	aggHandler, ok := aggregateRegistry[cmdType]
	if !ok {
		return nil, fmt.Errorf("invalid aggregate for: %s", cmdType)
	}

	agg, ok := aggHandler(cmd.AggregateID()).(cqrs.Aggregate)

	if !ok {
		return nil, fmt.Errorf("invalid aggregate for: %s", cmdType)
	}

	return agg, nil
}

// DecodeEvent decodes a raw event payload into its respective event type.
// It uses the event type to find the appropriate decoder and unmarshals the raw data.
// Returns an error if decoding fails or if no decoder is registered for the event type.
func DecodeEvent(evtType string, raw []byte) (cqrs.Event, error) {
	handler := eventDecoder[evtType]
	event, err := handler(raw)
	if err != nil {
		return nil, err
	}
	return event.(cqrs.Event), nil
}
