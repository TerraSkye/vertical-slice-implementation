package handlers

import (
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/cart/infrastructure"
)

// Automatically register the handler on package load
func init() {
	// Register event: CartCleared
	infrastructure.RegisterEvent(func(cart *domain.Cart) func(event *events.CartCleared) {
		return func(event *events.CartCleared) {
			cart.OnCartCleared(event)
		}
	})
}
