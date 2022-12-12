package api

import (
	"encoding/json"
	"net/http"

	"github.com/KeeseyLtd/microfun/users/internal/api/apierrors"
	"github.com/KeeseyLtd/microfun/users/internal/domain"
	"github.com/KeeseyLtd/microfun/users/internal/queries"
)

func LoginHandler(query queries.GetUserByLogin) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds domain.Credentials

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			errorResponse(w, apierrors.Response(apierrors.InvalidJson))

			return
		}

		if !creds.Validate() {
			errorResponse(w, apierrors.Response(apierrors.ValidationErr).Details(creds.Errors))

			return
		}

		login, _ := query.Handle(r.Context(), creds)

		successResponse(w, http.StatusOK, login)
	}
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

func RefreshHandler(query queries.GetUserByRefreshToken) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var token RefreshToken

		if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
			errorResponse(w, apierrors.Response(apierrors.InvalidJson))

			return
		}

		login, err := query.Handle(r.Context(), token.RefreshToken)

		if err != nil {
			errorResponse(w, apierrors.Response(apierrors.UnauthorizedErr))

			return
		}

		successResponse(w, http.StatusOK, login)
	}
}
