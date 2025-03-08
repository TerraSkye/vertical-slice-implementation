package archiveitem

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
	r.Methods("POST").Path("/api/commerce/carts/{id}/archive-item").Handler(
		kithttp.NewServer(
			func(ctx context.Context, request interface{}) (interface{}, error) {

				if err := s.ArchiveItem(ctx, request.(Payload)); err != nil {
					return nil, err
				}
				return struct{}{}, nil

			},
			decodeCreateRequest,
			infra.NoContent(),
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
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// decoding error no json provided
		return nil, err
	}

	return Payload{
		AggregateId: aggregateID,
		ProductId:   uuid.MustParse(payload.Data.ProductID),
	}, nil
}
