package handlers

import (
	"context"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
	"github.com/terraskye/vertical-slice-implementation/infra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"os"
)

func init() {
	cart.RegisterCommand(func(aggregate *domain.Cart) func(ctx context.Context, command *commands.AddItem) error {
		tracer := otel.Tracer("cart-service")

		return func(ctx context.Context, cmd *commands.AddItem) error {
			ctx, span := tracer.Start(ctx, "Cart:AddItem",
				trace.WithAttributes(
					// Add meta-related attributes
					attribute.String("cqrs.aggregate_id", cmd.AggregateId.String()),
					attribute.String("cqrs.aggregate_version", infra.MustExtractAggregateVersion(ctx)),
					attribute.String("cqrs.application", os.Getenv("application")),
					attribute.String("cqrs.causation_id", infra.MustExtractCausationId(ctx)),
					attribute.String("cqrs.correlation_id", trace.SpanContextFromContext(ctx).TraceID().String()),
					attribute.String("cqrs.command", "AddItem"),
					attribute.String("cqrs.function", "AddItem"),
					// Messaging attributes
					attribute.String("messaging.conversation_id", trace.SpanContextFromContext(ctx).TraceID().String()),
					attribute.String("messaging.destination", "AddItem"),
					attribute.String("messaging.destination_kind", "aggregate"),
					attribute.String("messaging.message_id", infra.MustExtractCausationId(ctx)),
					attribute.String("messaging.operation", "receive"),
					attribute.String("messaging.system", "cqrs"),
				),
			)
			defer span.End()
			err := aggregate.AddItem(ctx, cmd)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}

			return err
		}
	})
}
