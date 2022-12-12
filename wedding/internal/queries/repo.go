package queries

import (
	"context"

	"github.com/KeeseyLtd/microfun/wedding/internal/storage/postgres"
	"github.com/google/uuid"
)

type repository interface {
	GetWedding(ctx context.Context, id uuid.UUID) (postgres.Wedding, error)
	GetInvitations(context.Context, uuid.UUID) ([]postgres.Invitation, error)

	GetInvitation(ctx context.Context, arg postgres.GetInvitationParams) (postgres.Invitation, error)
}
