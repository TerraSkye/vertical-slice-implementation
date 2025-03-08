package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*ItemRemoved)(nil)

type ItemRemoved struct {
	AggregateId uuid.UUID
	ItemId      uuid.UUID
}

func (i ItemRemoved) AggregateID() uuid.UUID {
	return i.AggregateId
}
