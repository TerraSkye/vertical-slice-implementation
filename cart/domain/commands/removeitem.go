package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*RemoveItem)(nil)

type RemoveItem struct {
	AggregateId uuid.UUID
	ItemId      uuid.UUID
}

func (i RemoveItem) AggregateID() uuid.UUID {
	return i.AggregateId
}
