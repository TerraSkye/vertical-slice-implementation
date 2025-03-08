package infra

import (
	"context"
	"github.com/io-da/query"
)

type QueryHandler[T query.Query, R Readmodel] interface {
	HandleQuery(ctx context.Context, qry T) (R, error)
}

type QueryIteratorProvider interface {
	query.IteratorHandler
}
type QueryProvider interface {
	query.Handler
}
