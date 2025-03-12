package handlers

import (
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
)

func init() {
	cart.RegisterEvent(func(aggregate *domain.Pricing) func(event *events.PriceChanged) {
		return func(event *events.PriceChanged) {
			aggregate.OnPriceChanged(event)
		}
	})
}
