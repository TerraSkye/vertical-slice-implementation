package cartitems

import (
	"context"
	"github.com/google/uuid"
	"github.com/terraskye/vertical-slice-implementation/cqrs"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

// Query

type QueryHandler struct {
	store cqrs.EventStore
}

func NewQueryHandler(store cqrs.EventStore) infra.QueryHandler[Query, ReadModel] {
	return &QueryHandler{store: store}
}

func (q *QueryHandler) HandleQuery(ctx context.Context, qry Query) (ReadModel, error) {

	events, err := q.store.LoadFrom(ctx, qry.CartId, 0)

	if err != nil {
		return ReadModel{}, err
	}

	model := ReadModel{
		Items: make(map[uuid.UUID]*CartItem),
	}

	hydrate := infra.Hydrate(
		infra.NewHydrateHandler(model.OnItemRemoved),
		infra.NewHydrateHandler(model.OnCartCleared),
		infra.NewHydrateHandler(model.OnItemArchived),
		infra.NewHydrateHandler(model.OnCartCreated),
		infra.NewHydrateHandler(model.OnItemAdded),
	)

	for event := range events {
		hydrate(event)
	}

	return model, nil

}
