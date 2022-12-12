package queries

import (
	"context"

	"github.com/KeeseyLtd/microfun/wedding/internal/domain"
	"github.com/KeeseyLtd/microfun/wedding/internal/errs"
	"github.com/KeeseyLtd/microfun/wedding/internal/storage/postgres"
)

type GetInvitation struct {
	repo repository
}

func NewGetInvitationHandler(r repository) GetInvitation {
	return GetInvitation{repo: r}
}

func (h GetInvitation) Handle(ctx context.Context, args domain.GetInvitation) (domain.Invitation, error) {
	invite, err := h.repo.GetInvitation(ctx, postgres.GetInvitationParams{
		WeddingID: args.WeddingID,
		ID:        args.InvitationID,
	})

	if err != nil {
		return domain.Invitation{}, errs.ParseDBErrors(ctx, err, "invitation")
	}

	return domain.Invitation{
		ID:       invite.ID,
		Invitees: invite.Invitees,
		Status:   domain.InvitationStatus(invite.Status),
	}, nil
}
