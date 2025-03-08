package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*CartSubmitted)(nil)

type CartSubmitted struct {
	AggregateId     uuid.UUID
	OrderedProducts []interface{}
	TotalPrice      float64
}

func (c CartSubmitted) AggregateID() uuid.UUID {
	return c.AggregateId
}
