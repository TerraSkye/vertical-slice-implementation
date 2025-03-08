package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*PriceChanged)(nil)

type PriceChanged struct {
	NewPrice  float64
	OldPrice  float64
	ProductId uuid.UUID
}

func (p PriceChanged) AggregateID() uuid.UUID {
	return p.ProductId
}
