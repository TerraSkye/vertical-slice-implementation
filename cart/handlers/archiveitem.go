package handlers

import (
	"context"
	"fmt"
	"github.com/terraskye/vertical-slice-implementation/cart"
	"github.com/terraskye/vertical-slice-implementation/infra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"

	"github.com/terraskye/vertical-slice-implementation/cart/domain"
	"github.com/terraskye/vertical-slice-implementation/cart/domain/commands"
)

func init() {
	cart.RegisterCommand(func(aggregate *domain.Cart) func(ctx context.Context, cmd *commands.ArchiveItem) error {
		tracer := otel.Tracer("cart-service")

		return func(ctx context.Context, cmd *commands.ArchiveItem) error {
			ctx, span := tracer.Start(ctx, "Cart::ArchiveItem",
				trace.WithAttributes(
					// Add meta-related attributes
					attribute.String("cqrs.aggregate_id", cmd.AggregateId.String()),
					attribute.String("cqrs.aggregate_version", fmt.Sprintf("%d", infra.MustExtractAggregateVersion(ctx))),
					attribute.String("cqrs.application", os.Getenv("application")),
					attribute.String("cqrs.causation_id", infra.MustExtractCausationId(ctx)),
					attribute.String("cqrs.correlation_id", infra.MustExtractCorrelationId(ctx)),
					attribute.String("cqrs.command", "ArchiveItem"),
					attribute.String("cqrs.function", "ArchiveItem"),
					// Messaging attributes
					attribute.String("messaging.conversation_id", infra.MustExtractCorrelationId(ctx)),
					attribute.String("messaging.destination", "ArchiveItem"),
					attribute.String("messaging.destination_kind", "aggregate"),
					attribute.String("messaging.message_id", infra.MustExtractCausationId(ctx)),
					attribute.String("messaging.operation", "receive"),
					attribute.String("messaging.system", "cqrs"),
				),
			)
			defer span.End()
			err := aggregate.ArchiveItem(ctx, cmd)
			if err != nil {
				span.RecordError(err)
			} else {

			}

			return nil
		}
	})
}
