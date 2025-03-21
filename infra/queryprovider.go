package infra

import (
	"context"
	"fmt"
	"github.com/io-da/query"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
)

type GenericQueryHandler[T query.Query, R Readmodel] interface {
	HandleQuery(ctx context.Context, qry T) (R, error)
}

type QueryIteratorProvider interface {
	query.IteratorHandler
	RegisterHandler(handler GenericQueryHandler[query.Query, Readmodel])
}

type QueryProvider interface {
	query.Handler
	RegisterHandler(handler GenericQueryHandler[query.Query, Readmodel])
}

type handler struct {
	handlers map[string]GenericQueryHandler[query.Query, Readmodel]
}

func NewQueryHandler() QueryProvider {
	return &handler{
		handlers: make(map[string]GenericQueryHandler[query.Query, Readmodel]),
	}
}

func (t *handler) RegisterHandler(handler GenericQueryHandler[query.Query, Readmodel]) {
	var cmd query.Query
	queryType := cqrs.TypeName(cmd)
	// Store a type-erased function that preserves the correct signature
	if _, ok := t.handlers[queryType]; ok {
		panic("duplicate query handler" + queryType)
	}
	t.handlers[queryType] = handler
}

func (t *handler) Handle(ctx context.Context, qry query.Query, res *query.Result) error {
	cqrs.TypeName(qry)

	provider, exists := t.handlers[cqrs.TypeName(qry)]

	if !exists {
		return fmt.Errorf("unknown query type: %s", cqrs.TypeName(qry))
	}

	result, err := provider.HandleQuery(ctx, qry)

	if err != nil {
		return err
	}

	res.Add(result)
	res.Done()

	return nil
}

type iteratorHandler struct {
	handlers map[string]GenericQueryHandler[query.Query, Readmodel]
}

func NewQueryIteratorHandler() QueryIteratorProvider {
	return &iteratorHandler{
		handlers: make(map[string]GenericQueryHandler[query.Query, Readmodel]),
	}
}

func (t *iteratorHandler) RegisterHandler(handler GenericQueryHandler[query.Query, Readmodel]) {
	var cmd query.Query
	queryType := cqrs.TypeName(cmd)
	// Store a type-erased function that preserves the correct signature
	if _, ok := t.handlers[queryType]; ok {
		panic("duplicate query handler" + queryType)
	}
	t.handlers[queryType] = handler
}

func (t *iteratorHandler) Handle(ctx context.Context, qry query.Query, res *query.IteratorResult) error {
	cqrs.TypeName(qry)

	provider, exists := t.handlers[cqrs.TypeName(qry)]

	if !exists {
		return fmt.Errorf("unknown query type: %s", cqrs.TypeName(qry))
	}

	result, err := provider.HandleQuery(ctx, qry)

	if err != nil {
		return err
	}

	res.Yield(result)

	res.Done()

	return nil
}
