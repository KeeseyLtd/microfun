package apierrors

import (
	"net/http"
)

type ErrorCode int

const (
	WeddingResourceNotFound ErrorCode = 1000 + iota
	InvalidUUIDErr
	UnhandledDBError
	SomethingWentWrong
	InvalidJson
	ValidationErr
	UnauthorizedErr
	BadCookieErr
)

var errorMessage = map[ErrorCode]string{
	WeddingResourceNotFound: "wedding not found",
	InvalidUUIDErr:          "invalid UUID provided",
	UnhandledDBError:        "unhandled database error",
	SomethingWentWrong:      "something went wrong",
	InvalidJson:             "invalid json provided",
	ValidationErr:           "validation errors",
	UnauthorizedErr:         "invalid or no credentials provided",
	BadCookieErr:            "can not read cookie",
}

func (ec ErrorCode) Error() string {
	return errorMessage[ec]
}

var statusCode = map[ErrorCode]int{
	WeddingResourceNotFound: http.StatusNotFound,
	InvalidUUIDErr:          http.StatusBadRequest,
	UnhandledDBError:        http.StatusInternalServerError,
	SomethingWentWrong:      http.StatusInternalServerError,
	InvalidJson:             http.StatusBadRequest,
	ValidationErr:           http.StatusUnprocessableEntity,
	UnauthorizedErr:         http.StatusUnauthorized,
	BadCookieErr:            http.StatusBadRequest,
}

func (ec ErrorCode) HttpStatusCode() int {
	return statusCode[ec]
}

type APIErrorResponse struct {
	Message string            `json:"message"`
	Code    int               `json:"status"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e APIErrorResponse) Error() string {
	return e.Message
}

func (e APIErrorResponse) Details(details map[string]string) APIErrorResponse {
	return APIErrorResponse{
		Message: e.Message,
		Code:    e.Code,
		Errors:  details,
	}
}

func (e APIErrorResponse) SetCustomMessage(message string) APIErrorResponse {
	return APIErrorResponse{
		Message: message,
		Code:    e.Code,
	}
}

func APIResponse(e ErrorCode) APIErrorResponse {
	return APIErrorResponse{
		Message: e.Error(),
		Code:    e.HttpStatusCode(),
	}
}

func E(message string, code int) APIErrorResponse {
	return APIErrorResponse{
		Message: message,
		Code:    code,
	}
}
