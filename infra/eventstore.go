package infra

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"sync"
)

type MemoryStore struct {
	mu     sync.RWMutex
	bus    EventBus
	events map[uuid.UUID][]*cqrs.Envelope
}

func (m *MemoryStore) Save(ctx context.Context, events []cqrs.Envelope, originalVersion uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(events) == 0 {
		return nil
	}

	for i, event := range events {
		currentVersion := uint64(len(m.events[event.Event.AggregateID()]))

		if currentVersion != originalVersion {
			return fmt.Errorf("version did not match")
		}
		m.events[event.Event.AggregateID()] = append(m.events[event.Event.AggregateID()], &events[i])
		originalVersion++

	}

	for _, event := range events {
		//TODO populate the context with the event envelope
		if err := m.bus.Dispatch(ctx, event.Event); err != nil {
			//TODO now  we have an issue, where we recorded it happening but we cannot publish the events.
			return err
		}
	}

	return nil

}

func (m *MemoryStore) Load(ctx context.Context, u uuid.UUID) (<-chan *cqrs.Envelope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(chan *cqrs.Envelope, 10)
	go func() {
		defer close(out)
		if events, exists := m.events[u]; exists {
			for _, event := range events {
				select {
				case <-ctx.Done():
					return
				case out <- event:
				}
			}
		}
	}()

	return out, nil
}

func (m *MemoryStore) LoadFrom(ctx context.Context, id uuid.UUID, version int) (<-chan *cqrs.Envelope, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(chan *cqrs.Envelope, 10)
	go func() {
		defer close(out)
		if events, exists := m.events[id]; exists && version < len(events) {
			for _, event := range events[version:] {
				select {
				case <-ctx.Done():
					return
				case out <- event:
				}
			}
		}
	}()

	return out, nil
}

func (m *MemoryStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = make(map[uuid.UUID][]*cqrs.Envelope)
	return nil
}

func NewMemoryStore(bus EventBus) cqrs.EventStore {
	return &MemoryStore{
		events: make(map[uuid.UUID][]*cqrs.Envelope),
		bus:    bus,
	}
}
