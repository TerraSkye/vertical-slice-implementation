package events

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

// once it is installed, and loaded with the integration
var _ cqrs.Event = (*ItemAdded)(nil)

type ItemAdded struct {
	AggregateId       uuid.UUID
	Description       string
	Image             string
	Price             float64
	ItemId            uuid.UUID
	ProductId         uuid.UUID
	DeviceFingerPrint string
}

func (i ItemAdded) AggregateID() uuid.UUID {
	return i.AggregateId
}
