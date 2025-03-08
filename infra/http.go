package infra

import (
	"context"
	"net/http"
)

func NoContent() func(ctx context.Context, writer http.ResponseWriter, i any) error {
	return func(_ context.Context, writer http.ResponseWriter, _ any) error {
		writer.WriteHeader(http.StatusNoContent)
		return nil
	}
}
