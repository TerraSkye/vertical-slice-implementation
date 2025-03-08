package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*InventoryChanged)(nil)

type InventoryChanged struct {
	Inventory int
	ProductId uuid.UUID
}

func (i InventoryChanged) AggregateID() uuid.UUID {
	return i.ProductId
}
