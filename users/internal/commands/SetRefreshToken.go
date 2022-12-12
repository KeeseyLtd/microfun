package commands

import (
	"context"

	"github.com/KeeseyLtd/microfun/users/internal/domain"
	"github.com/KeeseyLtd/microfun/users/internal/logging"
	"github.com/google/uuid"
)

type SetRefreshToken struct{}

func NewSetRefreshTokenHandler() SetRefreshToken {
	return SetRefreshToken{}
}

func (h SetRefreshToken) Handle(ctx context.Context, refreshToken string, userID uuid.UUID) {
	domain.RefreshTokens.Store(refreshToken, userID)
	logging.WithContext(ctx).Info("added refresh token")
}
