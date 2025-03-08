package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*ArchiveItem)(nil)

type ArchiveItem struct {
	AggregateId uuid.UUID
	ProductId   uuid.UUID
}

func (i ArchiveItem) AggregateID() uuid.UUID {
	return i.AggregateId
}
