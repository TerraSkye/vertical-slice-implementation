package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*ChangePrice)(nil)

type ChangePrice struct {
	NewPrice  float64
	OldPrice  float64
	ProductId uuid.UUID
}

func (i ChangePrice) AggregateID() uuid.UUID {
	return i.ProductId
}
