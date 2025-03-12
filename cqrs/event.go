package cqrs

import (
	"github.com/google/uuid"
	"time"
)

type Envelope struct {
	UUID       uuid.UUID
	Metadata   map[string]any
	Event      Event
	Version    uint64
	OccurredAt time.Time
}

type Event interface {
	AggregateID() uuid.UUID
}
