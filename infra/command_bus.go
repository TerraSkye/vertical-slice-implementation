package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type CommandBus interface {
	Send(ctx context.Context, cmd cqrs.Command) <-chan error
	AddHandler(handler func(ctx context.Context, command cqrs.Command) error)
}
