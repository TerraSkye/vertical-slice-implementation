package domain

import (
	"context"
	"github.com/google/uuid"
	commands "github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	events "github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Pricing struct {
	*infra.AggregateBase
}

func (i *Pricing) New(uuid uuid.UUID) cqrs.Aggregate {
	return &Pricing{
		AggregateBase: infra.NewAggregateBase(uuid),
	}
}

func (p *Pricing) ChangePrice(ctx context.Context, cmd *commands.ChangePrice) error {
	p.AppendEvent(ctx, &events.PriceChanged{

		NewPrice:  cmd.NewPrice,
		OldPrice:  cmd.OldPrice,
		ProductId: cmd.ProductId,
	})
	return nil
}
func (p *Pricing) OnPriceChanged(event *events.PriceChanged) {}
