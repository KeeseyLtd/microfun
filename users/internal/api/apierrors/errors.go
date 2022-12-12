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

type Error struct {
	Message string            `json:"message"`
	Code    int               `json:"status"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func Response(e ErrorCode) Error {
	return Error{
		Message: e.Error(),
		Code:    e.HttpStatusCode(),
	}
}

func E(message string, code int) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}

func (e Error) Details(details map[string]string) Error {
	return Error{
		Message: e.Message,
		Code:    e.Code,
		Errors:  details,
	}
}
