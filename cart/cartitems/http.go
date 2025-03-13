package cartitems

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/io-da/query"
	"github.com/terraskye/vertical-slice-implementation/infra"
	"net/http"
)

func MakeHttpHandler(r *mux.Router, bus *query.Bus) http.Handler {
	queryHandler := infra.NewQueryGateway[*Query, ReadModel](bus)

	r.Methods("GET").Path("/api/commerce/carts/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/items").Handler(
		kithttp.NewServer(
			func(ctx context.Context, request interface{}) (interface{}, error) {
				model, err := queryHandler.Query(ctx, request.(*Query))

				if err != nil {
					return nil, err
				}

				return model.First(), nil

			},
			decodeCreateRequest,
			encodeResponse(),
		),
	)

	return r
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (any, error) {
	return &Query{CartId: uuid.MustParse(mux.Vars(r)["id"])}, nil
}

func encodeResponse() kithttp.EncodeResponseFunc {

	return func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {

		response := i.(*ReadModel)

		return json.NewEncoder(writer).Encode(response)

	}
}
