package commands

import (
	"context"

	"github.com/KeeseyLtd/microfun/wedding/internal/storage/postgres"
)

type repository interface {
	CreateWedding(ctx context.Context, arg postgres.CreateWeddingParams) (postgres.CreateWeddingRow, error)
}
