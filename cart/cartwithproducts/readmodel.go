package cartwithproducts

import (
	"github.com/google/uuid"
)

type CartsWithProductsReadModel struct {
	AggregateId uuid.UUID
	ProductId   uuid.UUID
}
