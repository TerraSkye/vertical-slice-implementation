package handlers

import (
	"context"
	domain "domain"
	commands "github.com/terraskye/vertical-slice-generator/gen/cart/domain/commands"
	infrastructure "github.com/terraskye/vertical-slice-implementation/cart/infrastructure"
)

func init() {
	infrastructure.RegisterCommand(func(aggregate *domain.Inventory) func(ctx context.Context, cmd *commands.ChangeInventory) error {
		return func(ctx context.Context, cmd *commands.ChangeInventory) error {
			return aggregate.ChangeInventory(ctx, cmd)
		}
	})
}
