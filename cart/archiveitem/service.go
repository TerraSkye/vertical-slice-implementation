package archiveitem

import (
	"context"
	"github.com/google/uuid"
	commands "github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

type Service interface {
	ArchiveItem(ctx context.Context, payload Payload) error
}

type Payload struct {
	AggregateId uuid.UUID
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
func (s *service) ArchiveItem(ctx context.Context, payload Payload) error {
	cmd := &commands.ArchiveItem{

		AggregateId: payload.AggregateId,
		ProductId:   payload.ProductId,
	}
	if err := <-s.commandBus.Send(ctx, cmd); err != nil {
		return err
	}

	return nil
}
