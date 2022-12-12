package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/KeeseyLtd/microfun/wedding/internal/domain"
	"github.com/KeeseyLtd/microfun/wedding/internal/errs"
)

type GetWedding struct {
	repo repository
}

func NewGetWeddingHandler(r repository) GetWedding {
	return GetWedding{repo: r}
}

func (h GetWedding) Handle(ctx context.Context, weddingID uuid.UUID) (domain.Wedding, error) {
	wedding, err := h.repo.GetWedding(ctx, weddingID)

	if err != nil {
		return domain.Wedding{}, errs.ParseDBErrors(ctx, err, "wedding")
	}

	dbInvites, err := h.repo.GetInvitations(ctx, weddingID)
	if err != nil {
		return domain.Wedding{}, errs.ParseDBErrors(ctx, err, "invitation")
	}

	invites := []domain.Invitation{}
	for _, i := range dbInvites {
		invites = append(invites, domain.Invitation{
			ID:       i.ID,
			Status:   domain.InvitationStatus(i.Status),
			Invitees: i.Invitees,
		})
	}

	weddingDTO := domain.Wedding{
		ID:          wedding.ID,
		Names:       wedding.Names,
		WeddingDate: wedding.WeddingDate,
		Invitations: invites,
	}

	return weddingDTO, nil
}
