package domain

import (
	"context"
	"github.com/google/uuid"
	commands "github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	events "github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Inventory struct {
	*infra.AggregateBase
}

func (i *Inventory) New(uuid uuid.UUID) cqrs.Aggregate {
	return &Inventory{
		AggregateBase: infra.NewAggregateBase(uuid),
	}
}

func (i *Inventory) ChangeInventory(ctx context.Context, cmd *commands.ChangeInventory) error {
	i.AppendEvent(ctx, &events.InventoryChanged{

		Inventory: cmd.Inventory,
		ProductId: cmd.ProductId,
	})
	return nil
}
func (i *Inventory) OnInventoryChanged(event *events.InventoryChanged) {}
