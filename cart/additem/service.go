package additem

import (
	"context"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Service interface {
	AddItem(ctx context.Context, payload Payload) error
}

type Payload struct {
	AggregateId uuid.UUID
	Description string
	Image       string
	ItemId      uuid.UUID
	ProductId   uuid.UUID
}
type service struct {
	commandBus infra.CommandBus
}

func NewService(commandBus infra.CommandBus) Service {
	return &service{
		commandBus: commandBus,
	}
}

func (s *service) AddItem(ctx context.Context, payload Payload) error {

	// price
	cmd := &commands.AddItem{

		AggregateId: payload.AggregateId,
		Description: payload.Description,
		Image:       payload.Image,
		ItemId:      payload.ItemId,
		Price:       995,
		ProductId:   payload.ProductId,
	}
	if err := <-s.commandBus.Send(ctx, cmd); err != nil {
		return err
	}

	return nil
}
