package handlers

import (
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
)

func init() {
	cart.RegisterEvent(func(aggregate *domain.Cart) func(event *events.ItemArchived) {
		return func(event *events.ItemArchived) {
			aggregate.OnItemArchived(event)
		}
	})
}
