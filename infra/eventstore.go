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
	events map[uuid.UUID][]cqrs.Event
}

func (m *MemoryStore) Save(ctx context.Context, events []cqrs.Event, originalVersion uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		currentVersion := uint64(len(m.events[event.AggregateID()]))

		if currentVersion != originalVersion {
			return fmt.Errorf("version did not match")
		}
		m.events[event.AggregateID()] = append(m.events[event.AggregateID()], event)

	}
	return nil

}

func (m *MemoryStore) Load(ctx context.Context, u uuid.UUID) (<-chan cqrs.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(chan cqrs.Event, 10)
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

func (m *MemoryStore) LoadFrom(ctx context.Context, id uuid.UUID, version int) (<-chan cqrs.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(chan cqrs.Event, 10)
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
	m.events = make(map[uuid.UUID][]cqrs.Event)
	return nil
}

func NewMemoryStore() cqrs.EventStore {
	return &MemoryStore{
		events: make(map[uuid.UUID][]cqrs.Event),
	}
}
