package additem

import (
	"context"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/terraskye/vertical-slice-implementation/infra"
	"net/http"
)

func MakeHttpHandler(r *mux.Router, s Service) http.Handler {
	r.Methods("POST").Path("/api/commerce/carts/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/add-item").Handler(
		kithttp.NewServer(
			func(ctx context.Context, request interface{}) (interface{}, error) { // 400-422
				if err := s.AddItem(ctx, request.(Payload)); err != nil {
					return nil, err
				}
				return struct{}{}, nil
			},
			decodeCreateRequest, // 400 - 404
			infra.NoContent(),
			kithttp.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
				// errors that can happen
			}),
		),
	)

	return r
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {

	aggregateID, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		return nil, fmt.Errorf("expected uuid but got %s", mux.Vars(r)["id"])
	}

	var payload struct {
		Data struct {
			ProductID string `json:"product_id" validate:"required,uuid4"`
			ItemID    string `json:"item_id" validate:"required,uuid4"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// decoding error no json provided
		return nil, err
	}

	return Payload{
		AggregateId: aggregateID,
		ProductId:   uuid.MustParse(payload.Data.ProductID),
		ItemId:      uuid.MustParse(payload.Data.ProductID),
	}, nil
}
