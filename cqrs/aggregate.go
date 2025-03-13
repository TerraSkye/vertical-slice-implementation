package cqrs

import (
	"context"
	"github.com/google/uuid"
)

// Aggregate is the interface that all aggregates must implement.
type Aggregate interface {

	// EntityID returns the unique identifier of the aggregate.
	EntityID() uuid.UUID

	// AggregateVersion returns the version of the aggregate.
	AggregateVersion() uint64

	// SetAggregateVersion sets the version of the aggregate.
	SetAggregateVersion(version uint64)

	// UncommittedEvents returns all the events that are currently uncommitted.
	UncommittedEvents() []Envelope

	// ClearUncommittedEvents clears all uncommitted events from the aggregate.
	ClearUncommittedEvents()

	// AppendEvent appends a new event to the aggregate's event list.
	AppendEvent(ctx context.Context, event Event, options ...EventOption)
}

type EventOption func(e *Envelope)

func WithMetaData(ctx context.Context) EventOption {
	return func(e *Envelope) {

		e.Metadata = map[string]any{
			//"user_id": claims.Subject,
			//"url":     afctx.GetCurrentUrl(ctx).String(),
		}
	}
}
