package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

type EventBus interface {
	Dispatch(ctx context.Context, event cqrs.Event) error
	Subscribe(handler EventHandler)
}

type eventBus struct {
	tracer   trace.Tracer
	handlers []EventHandler
	sync.RWMutex
}

func NewEventBus() EventBus {
	return &eventBus{}
}

// Dispatch sends the event to all subscribed handlers concurrently.
func (b *eventBus) Dispatch(ctx context.Context, event cqrs.Event) error {
	// Start a new tracing span, linking it to the incoming context
	ctx, span := b.tracer.Start(ctx, "EventBus.Dispatch",
		trace.WithAttributes(
			attribute.String("event.aggregate_id", event.AggregateID().String()),
			attribute.String("event.type", cqrs.TypeName(event)),
		),
	)
	defer span.End()

	b.RLock()
	handlers := append([]EventHandler{}, b.handlers...)
	b.RUnlock()

	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		expectedEvent := handler.NewEvent()

		if cqrs.TypeName(expectedEvent) == cqrs.TypeName(event) {
			wg.Add(1)
			go func(h EventHandler) {
				defer wg.Done()

				// Create a new span for the handler, linking it to Dispatch
				handlerCtx, handlerSpan := b.tracer.Start(ctx, "EventBus.HandleEvent",
					trace.WithAttributes(
						attribute.String("event.aggregate_id", event.AggregateID().String()),
						attribute.String("event.type", cqrs.TypeName(event)),
					),
				)

				if err := h.Handle(handlerCtx, event); err != nil {
					handlerSpan.RecordError(err)
					handlerSpan.SetStatus(codes.Error, err.Error())
					errChan <- err
				}

				handlerSpan.End()
			}(handler)
		}
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {

		err := <-errChan
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (b *eventBus) Subscribe(handler EventHandler) {
	b.Lock()
	defer b.Unlock()
	b.handlers = append(b.handlers, handler)
}
