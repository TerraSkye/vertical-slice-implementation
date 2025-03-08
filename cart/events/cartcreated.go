package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*CartCreated)(nil)

type CartCreated struct {
	AggregateId uuid.UUID
}

func (c CartCreated) AggregateID() uuid.UUID {
	return c.AggregateId
}
