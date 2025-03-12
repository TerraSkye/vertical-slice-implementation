package handlers

import (
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
)

func init() {
	cart.RegisterEvent(func(aggregate *domain.Cart) func(event *events.ItemAdded) {
		return func(event *events.ItemAdded) {
			aggregate.OnItemAdded(event)
		}
	})
}
