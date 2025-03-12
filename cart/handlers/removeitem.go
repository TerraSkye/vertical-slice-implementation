package handlers

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
)

func init() {
	cart.RegisterCommand(func(aggregate *domain.Cart) func(ctx context.Context, cmd *commands.RemoveItem) error {
		return func(ctx context.Context, cmd *commands.RemoveItem) error {
			return aggregate.RemoveItem(ctx, cmd)
		}
	})
}
