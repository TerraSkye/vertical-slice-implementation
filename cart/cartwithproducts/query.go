package cartwithproducts

import (
	"context"
	"github.com/google/uuid"
	"github.com/io-da/query"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

// Query
type Query struct {
	ProductID uuid.UUID
}

func (q Query) ID() []byte {
	return []byte(q.ProductID.String())
}

type QueryHandlerList struct{}

func (q *QueryHandlerList) HandleQuery(ctx context.Context, qry Query) (CartsWithProductsReadModel, error) {
	//TODO implement me
	panic("implement me")
}

func NewQueryHandler() infra.QueryHandler[Query, CartsWithProductsReadModel] {
	return &QueryHandlerList{}
}

func (q *QueryHandlerList) Handle(qry query.Query, res *query.IteratorResult) error {
	switch request := qry.(type) {
	case *Query:
		request.ProductID.Value()

		// fetch data

	}

	qry.ID()

	return nil
}
