package queries

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/KeeseyLtd/microfun/users/internal/commands"
	"github.com/KeeseyLtd/microfun/users/internal/config"
	"github.com/KeeseyLtd/microfun/users/internal/domain"
	"github.com/google/uuid"
)

type GetUserByRefreshToken struct {
	cfg *config.Config
}

func NewGetUserByRefreshTokenHandler(cfg *config.Config) GetUserByRefreshToken {
	return GetUserByRefreshToken{cfg: cfg}
}

func (h GetUserByRefreshToken) Handle(ctx context.Context, refreshToken string) (domain.SuccessfulLogin, error) {
	result, ok := domain.RefreshTokens.LoadAndDelete(refreshToken)
	if !ok {
		return domain.SuccessfulLogin{}, errors.New("refresh token not found")
	}

	user := domain.User{
		ID: result.(uuid.UUID),
	}

	jwt, _ := user.GenerateJWTToken(ctx, h.cfg.JWTSecret)

	newRefreshToken := sha256.Sum256([]byte(
		fmt.Sprintf("%d:%s", time.Now().Unix(), user.ID.String()),
	))

	newRefreshTokenString := fmt.Sprintf("%x", newRefreshToken)

	go commands.NewSetRefreshTokenHandler().Handle(context.Background(), newRefreshTokenString, user.ID)

	return domain.SuccessfulLogin{
		Token:        jwt,
		RefreshToken: newRefreshTokenString,
	}, nil
}
