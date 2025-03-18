package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CommandHandler struct {
	store  cqrs.EventStore
	tracer trace.Tracer
}

// NewCommandHandler creates a new CommandHandler for an aggregate type.
func NewCommandHandler(store cqrs.EventStore) *CommandHandler {
	h := &CommandHandler{
		store:  store,
		tracer: otel.Tracer("command-handler"),
	}

	return h
}

func (h *CommandHandler) Handle(ctx context.Context, command cqrs.Command) error {

	// Start the tracing span for handling the command
	ctx, span := h.tracer.Start(ctx, "HandleCommand",
		trace.WithAttributes(
			attribute.String("aggregate.id", command.AggregateID().String()),
			attribute.String("command.type", cqrs.TypeName(command)),
		),
	)
	defer span.End()

	aggregate, err := cart.AggregateForCommand(command)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	events, err := h.store.LoadFrom(ctx, command.AggregateID(), 0)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	var version uint64

	for event := range events {
		if err := cart.DispatchEvent(aggregate, event.Event); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return err
		}

		version++
	}
	aggregate.SetAggregateVersion(version)

	if err := cart.DispatchCommand(WithAggregateUUID(WithAggregateVersion(ctx, version), command.AggregateID()), aggregate, command); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	//events
	uncomittedEvents := aggregate.UncommittedEvents()

	if err := h.store.Save(ctx, uncomittedEvents, version); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	// Mark success in tracing
	span.SetStatus(codes.Ok, "")

	return nil

}
