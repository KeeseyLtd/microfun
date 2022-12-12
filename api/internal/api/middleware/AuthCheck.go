package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/apierrors"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/v1/handlers"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/logging" 
		"github.com/KeeseyLtd/microfun/api-gateway/internal/domain"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func AuthCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO This can be done better
		noAuth := []string{
			"/v1/auth/login",
			"/v1/auth/refresh",
		}

		for _, url := range noAuth {
			if url == r.URL.Path {
				h.ServeHTTP(w, r)

				return
			}
		}

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			handlers.ErrorResponse(w, apierrors.APIResponse(apierrors.UnauthorizedErr).SetCustomMessage("missing auth token"))

			return
		}

		splitToken := strings.Split(tokenString, " ")

		if len(splitToken) != 2 {
			handlers.ErrorResponse(w, apierrors.APIResponse(apierrors.UnauthorizedErr).SetCustomMessage("invalid auth token"))

			return
		}

		tokenPart := splitToken[1]
		var claims jwt.RegisteredClaims

		token, err := jwt.ParseWithClaims(tokenPart, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("TEST"), nil
		})

		if err != nil {
			var message string
			logging.WithContext(r.Context()).With("error", err).Info("error parsing jwt token with claims")
			if strings.Contains(err.Error(), "token is expired") {
				message = "token has expired. please login again"
			} else {
				message = "invalid auth token"
			}

			handlers.ErrorResponse(w, apierrors.APIResponse(apierrors.UnauthorizedErr).SetCustomMessage(message))

			return
		}

		if !token.Valid {
			handlers.ErrorResponse(w, apierrors.APIResponse(apierrors.UnauthorizedErr).SetCustomMessage("invalid auth token"))

			return
		}

		userID := uuid.MustParse(claims.Subject)

		ctx := logging.NewContext(r.Context(), "userId", userID)

		ctx = context.WithValue(ctx, domain.UserInfoKey, userID)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
