package handlers

import (
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

func init() {
	cart.RegisterAggregate(func(id uuid.UUID) *domain.Pricing {
		return &domain.Pricing{
			AggregateBase: infra.NewAggregateBase(id),
		}
	})
}
