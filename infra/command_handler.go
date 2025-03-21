package infra

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
)

type CommandHandler struct {
	store               cqrs.EventStore
	tracer              trace.Tracer
	aggregateForCommand func(cmd cqrs.Command) (cqrs.Aggregate, error)
	dispatchEvent       func(aggregate cqrs.Aggregate, event cqrs.Event) error
	dispatchCommand     func(ctx context.Context, aggregate cqrs.Aggregate, command cqrs.Command) error
}

// NewCommandHandler creates a new CommandHandler for an aggregate type.
func NewCommandHandler(
	store cqrs.EventStore,
	aggregateForCommand func(cmd cqrs.Command) (cqrs.Aggregate, error),
	dispatchEvent func(aggregate cqrs.Aggregate, event cqrs.Event) error,
	dispatchCommand func(ctx context.Context, aggregate cqrs.Aggregate, command cqrs.Command) error,
) *CommandHandler {

	tracer := otel.Tracer("command-handler")

	h := &CommandHandler{
		store:               store,
		tracer:              tracer,
		aggregateForCommand: aggregateForCommand,
		dispatchEvent: func(aggregate cqrs.Aggregate, event cqrs.Event) error {
			err := dispatchEvent(aggregate, event)
			return err
		},
		dispatchCommand: func(ctx context.Context, aggregate cqrs.Aggregate, command cqrs.Command) error {
			ctx, span := tracer.Start(ctx, "cqrs.command.handler.apply_command",
				trace.WithAttributes(
					attribute.String("command.type", cqrs.TypeName(command)),
					attribute.String("aggregate.id", command.AggregateID().String()),
					attribute.String("aggregate.version", MustExtractAggregateVersion(ctx)),

					attribute.String("cqrs.causation_id", MustExtractCausationId(ctx)),
					attribute.String("cqrs.correlation_id", trace.SpanContextFromContext(ctx).TraceID().String()),
					attribute.String("cqrs.command", cqrs.TypeName(command)),
					// Messaging attributes
					attribute.String("messaging.conversation_id", trace.SpanContextFromContext(ctx).TraceID().String()),
					attribute.String("messaging.destination_kind", "aggregate"),
					attribute.String("messaging.message_id", MustExtractCausationId(ctx)),
					semconv.MessagingOperationTypeProcess,
					attribute.String("messaging.system", "cqrs"),
				),
			)

			defer span.End()

			err := dispatchCommand(ctx, aggregate, command)

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")

			}
			return err

		},
	}

	return h
}

func (h *CommandHandler) Handle(ctx context.Context, command cqrs.Command) error {

	// Start the tracing span for handling the command
	ctx, span := h.tracer.Start(ctx, "cqrs.command.handler.execute",
		trace.WithAttributes(
			attribute.String("aggregate.id", command.AggregateID().String()),
			attribute.String("command.type", cqrs.TypeName(command)),
			semconv.MessagingOperationTypeReceive,
		),
	)
	defer span.End()

	aggregate, err := h.aggregateForCommand(command)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var version uint64

	{

		ctx2, span2 := h.tracer.Start(ctx, "cqrs.command.handler.load_aggregate",
			trace.WithAttributes(
				attribute.String("aggregate.id", command.AggregateID().String()),
			),
		)

		// the aggregate is being prepared here. we would want to trace this as wel.

		events, err := h.store.LoadFrom(ctx2, command.AggregateID(), 0)

		if err != nil {
			span2.RecordError(err)
			span2.SetStatus(codes.Error, err.Error())
			span2.End()
			return err
		}

		for event := range events {
			if err := h.dispatchEvent(aggregate, event.Event); err != nil {
				span2.RecordError(err)
				span2.SetStatus(codes.Error, err.Error())
				span2.End()
				return err
			}

			version++
		}
		aggregate.SetAggregateVersion(version)

		span2.AddEvent("aggregate loaded", trace.WithAttributes(
			attribute.String("aggregate.id", command.AggregateID().String()),
			attribute.Int64("aggregate.version", int64(version)),
		))

		span2.SetStatus(codes.Ok, "aggregate loaded successfully")
		span2.End()
	}

	if err := h.dispatchCommand(WithAggregateUUID(WithAggregateVersion(ctx, version), command.AggregateID()), aggregate, command); err != nil {
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
