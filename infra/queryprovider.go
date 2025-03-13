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

//
//type QueryIteratorProvider interface {
//	query.IteratorHandler
//}
//type QueryProvider interface {
//	query.Handler
//}
//
//type genericQueryHandler[T query.Query, R any] struct {
//}
//
//func NewQueryHandler[T query.Query, R any](bus *query.Bus) QueryHandler[T, R] {
//	return &genericQueryHandler[T, R]{}
//}
//
//func (r *genericQueryHandler[T, R]) Handle(ctx context.Context, query T) (R, error) {
//
//}

//func (h *GenericQueryGateway[T, R]) Query(ctx context.Context, qry T) (Res[R], error) {
//	result, err := h.bus.Query(qry)
//
//	var envelope Res[R]
//
//	if err != nil {
//		return envelope, err
//	}
//
//	if resultSize := len(result.All()); resultSize > 0 {
//
//		envelope = Res[R]{make([]*R, resultSize)}
//		for i, entity := range result.All() {
//			envelope.results[i] = entity.(*R)
//		}
//	}
//
//	return envelope, nil
//}

//func (h *GenericQueryGateway[T, R]) IteratorQuery(ctx context.Context, qry T) (<-chan R, error) {
//
//	result, err := h.bus.IteratorQuery(qry)
//
//	if err != nil {
//		return nil, err
//	}
//
//	output := make(chan R)
//
//	for res := range result.Iterate() {
//		output <- res
//	}
//
//	return output, nil
//}
