package cqrs

import (
	"context"
	"github.com/google/uuid"
)

// EventStore is an interface for an event sourcing event store.
type EventStore interface {
	// Save appends all events in the event stream to the store.
	Save(ctx context.Context, events []Event, originalVersion uint64) error

	// Load loads all events for the aggregate id from the store.
	Load(context.Context, uuid.UUID) (<-chan Event, error)

	// LoadFrom loads all events from version for the aggregate id from the store.
	LoadFrom(ctx context.Context, id uuid.UUID, version int) (<-chan Event, error)

	// Close closes the EventStore.
	Close() error
}
