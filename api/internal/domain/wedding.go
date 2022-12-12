package domain

import (
	"time"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/validation"
	"github.com/KeeseyLtd/microfun/api-gateway/pb"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateWeddingForm struct {
	// required: true
	// example: Daniella & Adam
	Names string `json:"names" validate:"required"`

	// required: true
	// minLength: 8
	// example: 2006-01-02 15:04:05
	WeddingDate string `json:"weddingDate" validate:"required,datetime=2006-01-02 15:04:05"`

	// swagger:ignore
	UserID uuid.UUID `json:"-"`

	// swagger:ignore
	Errors map[string]string `json:"errors,omitempty"`
}

func (dto *CreateWeddingForm) Validate() bool {
	dto.Errors = make(map[string]string)

	validate, trans := validation.GetValidator()

	if err := validate.Struct(dto); err != nil {
		errs := err.(validator.ValidationErrors)

		dto.Errors = validation.RemoveTopStruct(errs.Translate(trans))
	}

	return len(dto.Errors) == 0
}

type CreateWeddingParams struct {
	Names       string
	WeddingDate time.Time
	UserID      uuid.UUID
}

// swagger:model
type Wedding struct {
	// swagger:strfmt uuid
	ID uuid.UUID `json:"id"`

	// example: Dani and Adam
	Names string `json:"names"`

	// example: 2006-01-02 15:05:05
	WeddingDate time.Time    `json:"weddingDate"`
	Invitations []Invitation `json:"invitations"`
}

// swagger:model
type Invitation struct {
	// swagger:strfmt uuid
	ID uuid.UUID `json:"id"`

	// example: Mr & Mrs Jones
	Invitees string `json:"invitees"`

	// swagger:enum InvitationStatus
	Status InvitationStatus `json:"status"`
}

// swagger:enum InvitationStatus
type InvitationStatus int64

const (
	Pending      InvitationStatus = 0
	Attending    InvitationStatus = 1
	NotAttending InvitationStatus = 2
)

func (i InvitationStatus) String() string {
	var toString = map[InvitationStatus]string{
		Pending:      "Pending",
		Attending:    "Attending",
		NotAttending: "Not Attending",
	}

	return toString[i]
}

func NewWeddingFromProto(w *pb.Wedding) Wedding {
	wedding := Wedding{
		ID:          uuid.MustParse(w.Uuid),
		Names:       w.Names,
		WeddingDate: w.GetWeddingDate().AsTime(),
	}

	for _, invite := range w.Invitations {
		wedding.Invitations = append(wedding.Invitations, Invitation{
			ID:       uuid.MustParse(invite.Uuid),
			Invitees: invite.Names,
			Status:   InvitationStatus(invite.Status),
		})
	}

	return wedding
}
