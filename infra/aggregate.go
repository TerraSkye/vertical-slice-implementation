package infra

import (
	"context"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type AggregateBase struct {
	id     uuid.UUID
	v      uint64
	events []cqrs.Envelope
}

// NewAggregateBase creates an aggregate.
func NewAggregateBase(id uuid.UUID) *AggregateBase {
	return &AggregateBase{
		id:     id,
		events: make([]cqrs.Envelope, 0),
	}
}

// EntityID implements the EntityID method of the eh.Entity and eh.Aggregate interface.
func (a *AggregateBase) EntityID() uuid.UUID {
	return a.id
}

// AggregateVersion implements the AggregateVersion method of the Aggregate interface.
func (a *AggregateBase) AggregateVersion() uint64 {
	return a.v
}

// SetAggregateVersion implements the SetAggregateVersion method of the Aggregate interface.
func (a *AggregateBase) SetAggregateVersion(v uint64) {
	a.v = v
}

// UncommittedEvents implements the UncommittedEvents method of the eh.EventSource
// interface.
func (a *AggregateBase) UncommittedEvents() []cqrs.Envelope {
	return a.events
}

// ClearUncommittedEvents implements the ClearUncommittedEvents method of the eh.EventSource
// interface.
func (a *AggregateBase) ClearUncommittedEvents() {
	a.events = nil
}

// AppendEvent appends an event for later retrieval by Events().
func (a *AggregateBase) AppendEvent(ctx context.Context, event cqrs.Event, options ...cqrs.EventOption) {

	envelope := cqrs.Envelope{
		UUID:     uuid.New(),
		Metadata: make(map[string]any),
		Event:    event,
		Version:  a.AggregateVersion() + uint64(len(a.events)) + 1,
	}

	for _, option := range options {
		option(&envelope)
	}

	a.events = append(a.events, envelope)
}

type Readmodel interface {
}
