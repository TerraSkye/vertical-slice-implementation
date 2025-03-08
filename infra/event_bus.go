package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type EventBus interface {
	Dispatch(ctx context.Context, event cqrs.Event) error
	Subscribe(handler EventHandler)
}
