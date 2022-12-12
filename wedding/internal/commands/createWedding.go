package commands

import (
	"context"
	"fmt"

	"github.com/KeeseyLtd/microfun/wedding/internal/domain"
	"github.com/KeeseyLtd/microfun/wedding/internal/logging"
	"github.com/KeeseyLtd/microfun/wedding/internal/storage/postgres"
	"github.com/google/uuid"
)

type CreateWedding struct {
	repo repository
}

func NewCreateWeddingHandler(r repository) CreateWedding {
	return CreateWedding{repo: r}
}

func (h CreateWedding) Handle(ctx context.Context, args domain.CreateWeddingParams) error {
	weddingInsert := postgres.CreateWeddingParams{
		ID:          uuid.New(),
		Names:       args.Names,
		WeddingDate: args.WeddingDate,
		UserID:      args.UserID,
	}

	_, err := h.repo.CreateWedding(ctx, weddingInsert)
	if err != nil {
		return fmt.Errorf("unable to create wedding: %w", err)
	}

	logging.WithContext(ctx).With("newWeddingID", weddingInsert.ID).Info("created wedding")

	return nil
}
