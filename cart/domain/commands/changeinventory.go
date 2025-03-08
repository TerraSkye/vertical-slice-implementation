package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*ChangeInventory)(nil)

type ChangeInventory struct {
	Inventory int
	ProductId uuid.UUID
}

func (i ChangeInventory) AggregateID() uuid.UUID {
	return i.ProductId
}
