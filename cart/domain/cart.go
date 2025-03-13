package domain

import (
	"context"
	"errors"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Cart struct {
	*infra.AggregateBase
	submitted bool
	quantity  int
}

func (c *Cart) ClearCart(ctx context.Context, cmd *commands.ClearCart) error {
	c.AppendEvent(ctx, &events.CartCleared{
		AggregateId: cmd.AggregateId,
	})
	return nil
}
func (c *Cart) OnCartCleared(cmd *events.CartCleared) {
	c.quantity = 0
}

func (c *Cart) AddItem(ctx context.Context, cmd *commands.AddItem) error {

	if c.quantity == 10 {
		return errors.New("cannot add item to cart")
	}

	c.AppendEvent(ctx, &events.ItemAdded{
		AggregateId: cmd.AggregateId,
		Description: cmd.Description,
		Image:       cmd.Image,
		ItemId:      cmd.ItemId,
		Price:       cmd.Price,
		ProductId:   cmd.ProductId,
	})
	c.AppendEvent(ctx, &events.CartCreated{
		AggregateId: cmd.AggregateId,
	})
	return nil
}
func (c *Cart) OnItemAdded(cmd *events.ItemAdded) {
	c.quantity++
}

func (c *Cart) OnCartCreated(ev *events.CartCreated) {}

func (c *Cart) ArchiveItem(ctx context.Context, cmd *commands.ArchiveItem) error {
	c.AppendEvent(ctx, &events.ItemArchived{
		AggregateId: cmd.AggregateId,
		ItemId:      cmd.ProductId,
	})
	return nil
}
func (c *Cart) OnItemArchived(cmd *events.ItemArchived) {
	c.quantity--
}

func (c *Cart) SubmitCart(ctx context.Context, cmd *commands.SubmitCart) error {
	c.AppendEvent(ctx, &events.CartSubmitted{
		AggregateId:     cmd.AggregateId,
		OrderedProducts: cmd.OrderedProducts,
		//TotalPrice:      cmd.TotalPrice,
	})
	return nil
}
func (c *Cart) OnCartSubmitted(cmd *events.CartSubmitted) {

}

func (c *Cart) RemoveItem(ctx context.Context, cmd *commands.RemoveItem) error {
	c.AppendEvent(ctx, &events.ItemRemoved{
		AggregateId: cmd.AggregateId,
		ItemId:      cmd.ItemId,
	})
	return nil
}
func (c *Cart) OnItemRemoved(cmd *events.ItemRemoved) {
	c.quantity--
}
