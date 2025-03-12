package handlers

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
)

func init() {
	cart.RegisterCommand(func(aggregate *domain.Inventory) func(ctx context.Context, cmd *commands.ChangeInventory) error {
		return func(ctx context.Context, cmd *commands.ChangeInventory) error {
			return aggregate.ChangeInventory(ctx, cmd)
		}
	})
}
