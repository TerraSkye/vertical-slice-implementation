package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*AddItem)(nil)

type AddItem struct {
	AggregateId uuid.UUID
	Description string
	Image       string
	Price       float64
	ItemId      uuid.UUID
	ProductId   uuid.UUID
}

func (i AddItem) AggregateID() uuid.UUID {
	return i.AggregateId
}
