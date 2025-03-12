package handlers

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
)

func init() {
	cart.RegisterCommand(func(aggregate *domain.Cart) func(ctx context.Context, cmd *commands.ClearCart) error {
		return func(ctx context.Context, cmd *commands.ClearCart) error {
			return aggregate.ClearCart(ctx, cmd)
		}
	})
}
