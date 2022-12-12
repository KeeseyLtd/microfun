package queries

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/KeeseyLtd/microfun/users/internal/commands"
	"github.com/KeeseyLtd/microfun/users/internal/config"
	"github.com/KeeseyLtd/microfun/users/internal/domain"
	"github.com/google/uuid"
)

type GetUserByLogin struct {
	cfg *config.Config
}

func NewGetUserByLoginHandler(cfg *config.Config) GetUserByLogin {
	return GetUserByLogin{cfg: cfg}
}

func (h GetUserByLogin) Handle(ctx context.Context, args domain.Credentials) (domain.SuccessfulLogin, error) {
	user := domain.User{
		ID: uuid.New(),
	}

	jwt, _ := user.GenerateJWTToken(ctx, h.cfg.JWTSecret)

	refreshToken := sha256.Sum256([]byte(
		fmt.Sprintf("%d:%s", time.Now().Unix(), user.ID.String()),
	))

	refreshTokenString := fmt.Sprintf("%x", refreshToken)

	go commands.NewSetRefreshTokenHandler().Handle(context.Background(), refreshTokenString, user.ID)

	return domain.SuccessfulLogin{
		Token:        jwt,
		RefreshToken: refreshTokenString,
	}, nil
}
