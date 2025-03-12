package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"sync"
)

type EventBus interface {
	Dispatch(ctx context.Context, event cqrs.Event) error
	Subscribe(handler EventHandler)
}

type eventBus struct {
	handlers []EventHandler
	sync.RWMutex
}

func NewEventBus() EventBus {
	return &eventBus{}
}

// Dispatch sends the event to all subscribed handlers concurrently.
func (b *eventBus) Dispatch(ctx context.Context, event cqrs.Event) error {
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
				if err := h.Handle(ctx, event); err != nil {
					errChan <- err
				}
			}(handler)
		}
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (b *eventBus) Subscribe(handler EventHandler) {
	b.Lock()
	defer b.Unlock()
	b.handlers = append(b.handlers, handler)
}
