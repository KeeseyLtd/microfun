package domain

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/KeeseyLtd/microfun/users/internal/logging"
	"github.com/KeeseyLtd/microfun/users/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Credentials struct {
	Email string `json:"email" validate:"required,email"`

	Password string `json:"password" validate:"required"`

	Errors map[string]string `json:"errors,omitempty"`
}

func (dto *Credentials) Validate() bool {
	dto.Errors = make(map[string]string)

	validate, trans := validation.GetValidator()

	if err := validate.Struct(dto); err != nil {
		errs := err.(validator.ValidationErrors)

		dto.Errors = validation.RemoveTopStruct(errs.Translate(trans))
	}

	return len(dto.Errors) == 0
}

type User struct {
	// swagger:strfmt uuid
	ID uuid.UUID `json:"id"`

	// example: John
	FirstName string `json:"firstName"`

	// example: Smith
	LastName string `json:"lastName"`

	// swagger:strfmt email
	Email string `json:"email"`

	WeddingID uuid.UUID          `json:"-"`
	Status    VerificationStatus `json:"-"`
}

type VerificationStatus int

const (
	NotVerfied VerificationStatus = iota
	Verified
)

func (v VerificationStatus) String() string {
	var toString = map[VerificationStatus]string{
		NotVerfied: "Not Verified",
		Verified:   "Verified",
	}

	return toString[v]
}

func (u User) GenerateJWTToken(ctx context.Context, secret string) (string, error) {
	expTime := time.Now().Add(5 * time.Minute)

	claims := jwt.RegisteredClaims{
		Subject:   u.ID.String(),
		ExpiresAt: jwt.NewNumericDate(expTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logging.WithContext(ctx).Error(err)

		return "", errors.New("Error creating token")
	}

	return tokenString, nil
}

type SuccessfulLogin struct {
	// in:body
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

var RefreshTokens sync.Map
