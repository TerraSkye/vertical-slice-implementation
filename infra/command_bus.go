package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"sync"
)

type CommandBus interface {
	Send(ctx context.Context, cmd cqrs.Command) error
	AddHandler(handler func(ctx context.Context, command cqrs.Command) error)
}

type CommandWithCtx struct {
	Ctx        context.Context
	Command    cqrs.Command
	ResponseCh chan<- error
}

type commandBus struct {
	handlers []func(ctx context.Context, command cqrs.Command) error
	queue    chan CommandWithCtx
	sync.RWMutex
}

func NewCommandBus(bufferSize int) CommandBus {
	bus := &commandBus{
		queue: make(chan CommandWithCtx, bufferSize),
	}

	go bus.start()
	return bus
}

func (b *commandBus) Send(ctx context.Context, cmd cqrs.Command) error {
	responseCh := make(chan error, 1)

	// Enqueue the command with the response channel
	select {
	case b.queue <- CommandWithCtx{Ctx: ctx, Command: cmd, ResponseCh: responseCh}:
		// Wait for processing result
		select {
		case err := <-responseCh:
			return err // Return processing error (or nil if success)
		case <-ctx.Done():
			return ctx.Err() // Context timeout/cancellation
		}
	case <-ctx.Done():
		return ctx.Err() // Context timeout before enqueueing
	}
}

func (b *commandBus) AddHandler(handler func(ctx context.Context, command cqrs.Command) error) {
	b.Lock()
	defer b.Unlock()
	b.handlers = append(b.handlers, handler)
}

func (b *commandBus) start() {
	go func() {
		for cmdWithCtx := range b.queue {
			for _, handler := range b.handlers {

				go func(handlerFunc func(ctx context.Context, command cqrs.Command) error) {
					err := handlerFunc(cmdWithCtx.Ctx, cmdWithCtx.Command)
					if err != nil {
						cmdWithCtx.ResponseCh <- err
					} else {
						cmdWithCtx.ResponseCh <- nil
					}
				}(handler)
			}
		}
	}()
}
