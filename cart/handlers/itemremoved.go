package handlers

import (
	"context"
	domain "domain"
	events "github.com/terraskye/vertical-slice-generator/gen/cart/events"
	infrastructure "github.com/terraskye/vertical-slice-implementation/cart/infrastructure"
)

func init() {
	infrastructure.RegisterCommand(func(aggregate *domain.Cart) func(ctx context.Context, cmd *events.ItemRemoved) error {
		return func(ctx context.Context, cmd *events.ItemRemoved) error {
			return aggregate.OnItemRemoved(ctx, cmd)
		}
	})
}
