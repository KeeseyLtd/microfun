package errs

import (
	"context"
	"strings"

	"github.com/KeeseyLtd/microfun/wedding/internal/logging"
)

type ErrorCode int

const (
	WeddingResourceNotFound ErrorCode = 1000 + iota
	InvalidUUIDErr
	UnhandledDBError
	InviteResourceNotFound
	InvalidJson
	ValidationErr
)

var errorMessage = map[ErrorCode]string{
	WeddingResourceNotFound: "wedding not found",
	InvalidUUIDErr:          "invalid UUID provided",
	UnhandledDBError:        "unhandled database error",
	InviteResourceNotFound:  "invitation not found",
	InvalidJson:             "invalid json provided",
	ValidationErr:           "validation errors",
}

func (ec ErrorCode) Error() string {
	return errorMessage[ec]
}

func ParseDBErrors(ctx context.Context, err error, resource string) ErrorCode {
	if strings.Contains(err.Error(), "no rows in result set") {
		switch resource {
		case "wedding":
			return WeddingResourceNotFound
		case "invitation":
			return InviteResourceNotFound
		}
	}

	logging.WithContext(ctx).With("error", err).Error(UnhandledDBError)

	return UnhandledDBError
}
