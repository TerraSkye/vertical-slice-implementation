package commands

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Command = (*ClearCart)(nil)

type ClearCart struct {
	AggregateId uuid.UUID
}

func (i ClearCart) AggregateID() uuid.UUID {
	return i.AggregateId
}
