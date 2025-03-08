package cartwithproducts

import (
	"context"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cart/events"
)

type Projecter interface {
	OnItemArchived(ctx context.Context, ev *events.ItemArchived) error
	OnItemRemoved(ctx context.Context, ev *events.ItemRemoved) error
	OnItemAdded(ctx context.Context, ev *events.ItemAdded) error
	OnCartCreated(ctx context.Context, ev *events.CartCreated) error
	OnCartCleared(ctx context.Context, ev *events.CartCleared) error
}

type projector struct {
	entities map[uuid.UUID]*CartsWithProductsReadModel
}

func NewProjector() Projecter {
	return &projector{}
}

func (p *projector) OnItemArchived(ctx context.Context, ev *events.ItemArchived) error {

	return nil
}

func (p *projector) OnItemRemoved(ctx context.Context, ev *events.ItemRemoved) error {
	return nil
}

func (p *projector) OnItemAdded(ctx context.Context, ev *events.ItemAdded) error {
	return nil
}

func (p *projector) OnCartCreated(ctx context.Context, ev *events.CartCreated) error {
	return nil
}

func (p *projector) OnCartCleared(ctx context.Context, ev *events.CartCleared) error {
	return nil
}
