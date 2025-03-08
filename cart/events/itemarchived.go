package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*ItemArchived)(nil)

type ItemArchived struct {
	AggregateId uuid.UUID
	ItemId      uuid.UUID
}

func (i ItemArchived) AggregateID() uuid.UUID {
	return i.AggregateId
}
