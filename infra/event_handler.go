package infra

import (
	"context"
)

type EventHandler interface {
	// HandlerName is the name used in message.Router while creating handler.
	//
	// It will be also passed to EventsSubscriberConstructor.
	// May be useful, for example, to create a consumer group per each handler.
	//
	// WARNING: If HandlerName was changed and is used for generating consumer groups,
	// it may result with **reconsuming all messages** !!!
	HandlerName() string
	NewEvent() any
	Handle(ctx context.Context, event any) error
}

type genericEventHandler[T any] struct {
	handleFunc  func(ctx context.Context, event *T) error
	handlerName string
}

// NewEventHandler creates a new EventHandler implementation based on provided function
// and event type inferred from function argument.
func NewEventHandler[T any](
	handlerName string,
	handleFunc func(ctx context.Context, event *T) error,
) EventHandler {
	return &genericEventHandler[T]{
		handleFunc:  handleFunc,
		handlerName: handlerName,
	}
}

func (c genericEventHandler[T]) HandlerName() string {
	return c.handlerName
}

func (c genericEventHandler[T]) NewEvent() any {
	tVar := new(T)
	return tVar
}

func (c genericEventHandler[T]) Handle(ctx context.Context, e any) error {
	event := e.(*T)
	return c.handleFunc(ctx, event)
}

type GroupEventHandler interface {
	NewEvent() any
	Handle(ctx context.Context, event any) error
}

// NewGroupEventHandler creates a new GroupEventHandler implementation based on provided function
// and event type inferred from function argument.
func NewGroupEventHandler[T any](handleFunc func(ctx context.Context, event *T) error) GroupEventHandler {
	return &genericEventHandler[T]{
		handleFunc: handleFunc,
	}
}
