package handlers

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/cart/infrastructure"
)

func init() {
	infrastructure.RegisterCommand(func(cart *domain.Cart) func(ctx context.Context, cmd *commands.ClearCart) error {
		return func(ctx context.Context, cmd *commands.ClearCart) error {
			return cart.ClearCart(ctx, cmd)
		}
	})
}
