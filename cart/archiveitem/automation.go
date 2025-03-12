package archiveitem

import (
	"context"
	"github.com/io-da/query"
	"github.com/terraskye/vertical-slice-implementation/cart/cartwithproducts"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Automation interface {
	OnPriceChanged(ctx context.Context, ev *events.PriceChanged) (err error)
}
type automation struct {
	commandBus            infra.CommandBus
	cartItemsWithProducts infra.GenericQueryGateway[cartwithproducts.Query, cartwithproducts.CartsWithProductsReadModel]
}

func NewAutomation(commandBus infra.CommandBus, bus *query.Bus) Automation {
	return &automation{
		commandBus:            commandBus,
		cartItemsWithProducts: infra.NewQueryGateway[cartwithproducts.Query, cartwithproducts.CartsWithProductsReadModel](bus),
	}
}
func (a *automation) OnPriceChanged(ctx context.Context, ev *events.PriceChanged) error {

	result, err := a.cartItemsWithProducts.IteratorQuery(ctx, cartwithproducts.Query{ProductID: ev.ProductId})

	if err != nil {
		return err
	}

	for cartWithProduct := range result {
		cmd := &commands.ArchiveItem{
			AggregateId: cartWithProduct.AggregateId,
			ProductId:   ev.ProductId,
		}
		if err := a.commandBus.Send(ctx, cmd); err != nil {
			return err
		}

	}

	return nil
}
