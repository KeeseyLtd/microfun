package middleware

import (
	"net/http"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/logging"
	"github.com/google/uuid"
)

type ContextKey string

const RequestKey ContextKey = "RequestId"

func AddRequestUUID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.NewContext(r.Context(), "requestId", uuid.New().String())

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
