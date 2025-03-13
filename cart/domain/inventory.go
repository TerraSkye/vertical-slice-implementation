package domain

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Inventory struct {
	*infra.AggregateBase
}

func (i *Inventory) ChangeInventory(ctx context.Context, cmd *commands.ChangeInventory) error {
	i.AppendEvent(ctx, &events.InventoryChanged{

		Inventory: cmd.Inventory,
		ProductId: cmd.ProductId,
	})
	return nil
}
func (i *Inventory) OnInventoryChanged(event *events.InventoryChanged) {}
