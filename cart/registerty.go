package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// Generic command and event handler types
type CommandHandler[A any, T any] func(aggregate A) func(ctx context.Context, command T) error
type EventHandler[A any, T cqrs.Event] func(aggregate A) func(event T)
type AggregateHandler[A any] func(id uuid.UUID) A

var (
	commandRegistry    = make(map[string]any)
	eventRegistry      = make(map[string]any)
	aggregateRegistry  = make(map[string]func(id uuid.UUID) any)
	commandToAggregate = make(map[string]string)
	eventDecoder       = make(map[string]func(raw []byte) (any, error))
)

func AggregateForCommand(cmd cqrs.Command) (cqrs.Aggregate, error) {
	cmdType := cqrs.TypeName(cmd)

	aggType, ok := commandToAggregate[cmdType]

	if !ok {
		return nil, fmt.Errorf("no command to aggregate mapping found for: %s", cmdType)
	}

	aggHandler, ok := aggregateRegistry[aggType]
	if !ok {
		return nil, fmt.Errorf("invalid aggregate for: %s", cmdType)
	}

	agg, ok := aggHandler(cmd.AggregateID()).(cqrs.Aggregate)

	if !ok {
		return nil, fmt.Errorf("invalid aggregate for: %s", cmdType)
	}

	return agg, nil
}

// RegisterCommand registers a command handler
func RegisterCommand[A cqrs.Aggregate, T cqrs.Command](handler CommandHandler[A, T]) {
	var cmd T
	cmdType := cqrs.TypeName(cmd)

	// Store a type-erased function that preserves the correct signature
	commandRegistry[cmdType] = func(aggregate cqrs.Aggregate) func(context.Context, cqrs.Command) error {
		return func(ctx context.Context, command cqrs.Command) error {
			return handler(aggregate.(A))(ctx, command.(T)) // Explicit type assertion
		}
	}

	var agg A
	aggType := cqrs.TypeName(agg)
	commandToAggregate[cmdType] = aggType
}

// DispatchCommand finds and executes a command handler
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

// DispatchEvent finds and executes an event handler
func DispatchEvent[A cqrs.Aggregate, T cqrs.Event](aggregate A, event T) error {
	evtType := cqrs.TypeName(event)
	handlerRaw, exists := eventRegistry[evtType]
	if !exists {
		return errors.New("no handler registered for event: " + evtType)
	}

	// Ensure type safety
	handlerWrapper, ok := handlerRaw.(func(cqrs.Aggregate) func(cqrs.Event))
	if !ok {
		return errors.New("invalid event handler type for: " + evtType)
	}

	// Execute the command handler
	handlerWrapper(aggregate)(event)
	return nil
}

// RegisterEvent registers an event handler
func RegisterEvent[A cqrs.Aggregate, T cqrs.Event](handler EventHandler[A, T]) {

	var evt T
	evtType := cqrs.TypeName(evt)
	eventRegistry[evtType] = func(aggregate cqrs.Aggregate) func(cqrs.Event) {
		return func(event cqrs.Event) {
			handler(aggregate.(A))(event.(T)) // Explicit type assertion
		}
	}

	// Store an event decoder
	eventDecoder[evtType] = func(raw []byte) (any, error) {
		var evt T
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return evt, nil
	}
}

func RegisterAggregate[A cqrs.Aggregate](handler AggregateHandler[A]) {
	var evt A
	aggType := cqrs.TypeName(evt)
	aggregateRegistry[aggType] = func(id uuid.UUID) any {
		return handler(id)
	}
}

func DecodeEvent(evtType string, raw []byte) (cqrs.Event, error) {
	handler := eventDecoder[evtType]
	event, err := handler(raw)
	if err != nil {
		return nil, err
	}
	return event.(cqrs.Event), nil
}
