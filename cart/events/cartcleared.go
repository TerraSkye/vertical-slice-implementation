package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*CartCleared)(nil)

type CartCleared struct {
	AggregateId uuid.UUID
}

func (c CartCleared) AggregateID() uuid.UUID {
	return c.AggregateId
}
