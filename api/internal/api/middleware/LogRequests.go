package middleware

import (
	"net/http"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/logging"
)

func LogRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.WithContext(r.Context()).Info(r.URL)

		h.ServeHTTP(w, r)
	})
}
