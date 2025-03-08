package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*SubmitCart)(nil)

type SubmitCart struct {
	AggregateId     uuid.UUID
	OrderedProducts []interface{}
}

func (i SubmitCart) AggregateID() uuid.UUID {
	return i.AggregateId
}
