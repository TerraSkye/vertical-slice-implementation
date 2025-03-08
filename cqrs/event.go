package cqrs

import "github.com/google/uuid"

type Event interface {
	AggregateID() uuid.UUID
}
