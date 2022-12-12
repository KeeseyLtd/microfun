package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KeeseyLtd/microfun/users/internal/api/apierrors"
	"github.com/KeeseyLtd/microfun/users/internal/app"
	"github.com/go-chi/chi/v5"
)

func Routes(
	r *chi.Mux,
	a app.App,
) {
	r.Post("/login", LoginHandler(a.Queries.GetUserByLogin))
	r.Post("/refresh", RefreshHandler(a.Queries.GetUserByRefreshToken))

}

func errorResponse(w http.ResponseWriter, err error) {

	var apiError apierrors.Error
	if errors.As(err, &apiError) {
		w.WriteHeader(apiError.Code)

		json.NewEncoder(w).Encode(apiError)

		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func successResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	if data != nil {
		// a := data.(responsei)
		// a.Success(w)
		json.NewEncoder(w).Encode(data)
	}
}
