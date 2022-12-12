package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/apierrors"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/domain"
)

func ErrorResponse(w http.ResponseWriter, err error) {

	var apiError apierrors.APIErrorResponse
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

		json.NewEncoder(w).Encode(data)
	}
}

// swagger:response WeddingResponse
type WeddingResponse struct {
	// in: body
	Body struct {
		Data domain.Wedding `json:"data"`
	}
}

func (r WeddingResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Body)
}

// swagger:response LoginResponse
type LoginResponse struct {
	// in: body
	Body struct {
		Data SuccessfulLogin `json:"data"`
	}
}

func (r LoginResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Body)
}
