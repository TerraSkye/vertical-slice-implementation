package infrastructure

import (
	"context"
	"errors"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// Generic command and event handler types
type commandHandler[A any, T cqrs.Command] func(aggregate A) func(ctx context.Context, command T) error
type eventHandler[A any, T cqrs.Event] func(aggregate A) func(event T)

// Command and event registries
var commandRegistry = make(map[string]any)
var eventRegistry = make(map[string]any)

var aggregateRegistry = make(map[string]any)

// RegisterCommand registers a command handler
func RegisterCommand[A any, T cqrs.Command](handler commandHandler[A, T]) {
	var cmd T // Get the command type
	cmdType := cqrs.TypeName(cmd)
	commandRegistry[cmdType] = handler
	var aggregate A
	aggregateRegistry[cqrs.TypeName(aggregate)] = aggregate
}

// RegisterEvent registers an event handler
func RegisterEvent[A any, T cqrs.Event](handler eventHandler[A, T]) {
	var evt T // Get the event type
	evtType := cqrs.TypeName(evt)
	eventRegistry[evtType] = handler
}

// DispatchCommand finds and executes a command handler
func DispatchCommand[A any, T cqrs.Command](ctx context.Context, aggregate A, command T) error {
	cmdType := cqrs.TypeName(command)
	handlerRaw, exists := commandRegistry[cmdType]
	if !exists {
		return errors.New("no handler registered for command: " + cmdType)
	}

	// Type assertion to expected function signature
	handler, ok := handlerRaw.(commandHandler[A, T])
	if !ok {
		return errors.New("invalid command handler type for: " + cmdType)
	}

	// Execute the command
	return handler(aggregate)(ctx, command)
}

// DispatchEvent finds and executes an event handler
func DispatchEvent[A any, T cqrs.Event](aggregate A, event T) error {
	evtType := cqrs.TypeName(event)
	handlerRaw, exists := eventRegistry[evtType]
	if !exists {
		return errors.New("no handler registered for event: " + evtType)
	}

	// Type assertion to expected function signature
	handler, ok := handlerRaw.(eventHandler[A, T])
	if !ok {
		return errors.New("invalid event handler type for: " + evtType)
	}

	// Execute the event
	handler(aggregate)(event)
	return nil
}
