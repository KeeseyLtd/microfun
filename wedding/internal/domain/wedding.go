package domain

import (
	"time"

	"github.com/KeeseyLtd/microfun/wedding/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Wedding struct {
	ID          uuid.UUID    `json:"id"`
	Names       string       `json:"names"`
	WeddingDate time.Time    `json:"weddingDate"`
	Invitations []Invitation `json:"invitations"`
}

type Invitation struct {
	ID       uuid.UUID        `json:"id"`
	Invitees string           `json:"invitees"`
	Status   InvitationStatus `json:"status"`
}

type InvitationStatus int64

const (
	Pending      InvitationStatus = 0
	Attending    InvitationStatus = 1
	NotAttending InvitationStatus = 2
)

func (i InvitationStatus) String() string {
	return map[InvitationStatus]string{
		Pending:      "Pending",
		Attending:    "Attending",
		NotAttending: "Not Attending",
	}[i]
}

type GetInvitation struct {
	WeddingID    uuid.UUID
	InvitationID uuid.UUID
}

type CreateWeddingParams struct {
	Names       string
	WeddingDate time.Time
	UserID      uuid.UUID
}

func (w Wedding) ToProto() *pb.Wedding {
	wedding := pb.Wedding{
		Uuid:        w.ID.String(),
		Names:       w.Names,
		WeddingDate: timestamppb.New(w.WeddingDate),
	}

	for _, invite := range w.Invitations {
		wedding.Invitations = append(wedding.Invitations, &pb.Invitation{
			Uuid:   invite.ID.String(),
			Names:  invite.Invitees,
			Status: pb.Invitation_Status(invite.Status),
		})
	}

	return &wedding
}
