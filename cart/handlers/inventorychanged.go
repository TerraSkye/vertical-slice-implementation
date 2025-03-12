package handlers

import (
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
)

func init() {
	cart.RegisterEvent(func(aggregate *domain.Inventory) func(event *events.InventoryChanged) {
		return func(event *events.InventoryChanged) {
			aggregate.OnInventoryChanged(event)
		}
	})
}
