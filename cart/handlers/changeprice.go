package handlers

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
)

func init() {
	cart.RegisterCommand(func(aggregate *domain.Pricing) func(ctx context.Context, cmd *commands.ChangePrice) error {
		return func(ctx context.Context, cmd *commands.ChangePrice) error {
			return aggregate.ChangePrice(ctx, cmd)
		}
	})
}
