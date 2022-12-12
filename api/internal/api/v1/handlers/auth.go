package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/apierrors"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/logging"
)

type SuccessfulLogin struct {
	// in:body
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("POST", "http://user-service:9001/login", r.Body)
	if err != nil {
		logging.WithContext(r.Context()).With("error", err).Error("failed to create request")

		ErrorResponse(w, apierrors.APIResponse(apierrors.SomethingWentWrong))

		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.WithContext(r.Context()).With("error", err).Error("failed to make request")

		ErrorResponse(w, apierrors.APIResponse(apierrors.SomethingWentWrong))

		return
	}

	var loginResult SuccessfulLogin
	if err := json.NewDecoder(res.Body).Decode(&loginResult); err != nil {
		ErrorResponse(w, apierrors.APIResponse(apierrors.InvalidJson))

		return
	}

	lr := LoginResponse{}
	lr.Body.Data = loginResult

	http.SetCookie(w, &http.Cookie{
		Name:     "Refresh-Token",
		Value:    loginResult.RefreshToken,
		HttpOnly: true,
		Path:     "/v1/auth/refresh",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	successResponse(w, 200, lr)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("Refresh-Token")

	body, _ := json.Marshal(map[string]string{"refreshToken": token.Value})

	if err != nil {
		ErrorResponse(w, apierrors.APIResponse(apierrors.BadCookieErr))

		return
	}

	req, err := http.NewRequest("POST", "http://user-service:9001/refresh", bytes.NewBuffer(body))
	if err != nil {
		logging.WithContext(r.Context()).With("error", err).Error("failed to create request")

		ErrorResponse(w, apierrors.APIResponse(apierrors.SomethingWentWrong))

		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.WithContext(r.Context()).With("error", err).Error("failed to make request")

		ErrorResponse(w, apierrors.APIResponse(apierrors.SomethingWentWrong))

		return
	}

	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)

		b, _ := ioutil.ReadAll(res.Body)
		w.Write(b)

		return
	}

	var loginResult SuccessfulLogin
	if err := json.NewDecoder(res.Body).Decode(&loginResult); err != nil {
		ErrorResponse(w, apierrors.APIResponse(apierrors.InvalidJson))

		return
	}

	lr := LoginResponse{}
	lr.Body.Data = loginResult

	http.SetCookie(w, &http.Cookie{
		Name:     "Refresh-Token",
		Value:    loginResult.RefreshToken,
		HttpOnly: true,
		Path:     "/v1/auth/refresh",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	successResponse(w, 200, lr)
}
