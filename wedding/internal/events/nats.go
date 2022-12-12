package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/KeeseyLtd/microfun/wedding/internal/domain"
	"github.com/KeeseyLtd/microfun/wedding/internal/errs"
	"github.com/KeeseyLtd/microfun/wedding/internal/logging"
	"github.com/KeeseyLtd/microfun/wedding/pb"
	"github.com/google/uuid"
	nats "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	GetWeddingTopic    = "wedding.get"
	GetInvitationTopic = "wedding.invitation.get"

	CreateWeddingTopic = "wedding.create"
)

type MessageType int64

const (
	Success MessageType = 0
	Error   MessageType = 1
)

type GetWeddingResponseMessage struct {
	Type      MessageType
	Payload   domain.Wedding
	ErrorCode errs.ErrorCode
}

type GetInvitationResponseMessage struct {
	Type      MessageType
	Payload   domain.Invitation
	ErrorCode errs.ErrorCode
}

type Nats struct {
	nc *nats.Conn
}

func New(url string) (*Nats, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Nats{nc: nc}, nil
}

type getInvitationQueryer interface {
	Handle(ctx context.Context, args domain.GetInvitation) (domain.Invitation, error)
}

func (n *Nats) SubscribeGetInvitation(query getInvitationQueryer) {
	n.nc.Subscribe(GetInvitationTopic, func(msg *nats.Msg) {
		ctx := context.Background()
		m := GetInvitationResponseMessage{}
		args := domain.GetInvitation{}

		buf := bytes.NewBuffer(msg.Data)
		if err := gob.NewDecoder(buf).Decode(&args); err != nil {
			logging.WithContext(ctx).With("error", err).Error("could not decode message")
			return
		}

		invite, err := query.Handle(ctx, args)
		if err != nil {
			m.Type = Error
			m.ErrorCode = err.(errs.ErrorCode)
			respondGob(ctx, msg, m)

			return
		}

		m.Type = Success
		m.Payload = invite

		logging.WithContext(ctx).With("invitationID", invite.ID).Info("responding with invitation")

		respondGob(ctx, msg, m)
	})
}

type createWeddingCommand interface {
	Handle(ctx context.Context, args domain.CreateWeddingParams) error
}

func (n *Nats) SubscribeCreateWedding(command createWeddingCommand) {
	n.nc.Subscribe(CreateWeddingTopic, func(msg *nats.Msg) {
		ctx := context.Background()

		args := domain.CreateWeddingParams{}

		buf := bytes.NewBuffer(msg.Data)
		if err := gob.NewDecoder(buf).Decode(&args); err != nil {
			logging.WithContext(ctx).With("error", err).Error("could not decode message")
			return
		}

		err := command.Handle(ctx, args)

		if err != nil {
			logging.WithContext(ctx).With("error", err).Error("could not create wedding")
		}
	})
}

type getWeddingQueryer interface {
	Handle(ctx context.Context, weddingID uuid.UUID) (domain.Wedding, error)
}

func (n *Nats) SubscribeGetWedding(query getWeddingQueryer) {
	n.nc.Subscribe(GetWeddingTopic, func(msg *nats.Msg) {
		ctx := context.Background()
		m := &pb.WeddingResponse{}

		id := pb.GetWedding{}
		proto.Unmarshal(msg.Data, &id)

		parsedID, err := uuid.Parse(id.GetUuid())
		if err != nil {
			m.Type = pb.WeddingResponse_Error
			m.ErrorCode = int32(errs.InvalidUUIDErr)
			logging.WithContext(ctx).Warn(err)

			bytesSlice, _ := proto.Marshal(m)

			msg.Respond(bytesSlice)
			return
		}

		wedding, err := query.Handle(ctx, parsedID)
		if err != nil {
			m.Type = pb.WeddingResponse_Error
			m.ErrorCode = int32(err.(errs.ErrorCode))

			bytesSlice, _ := proto.Marshal(m)

			msg.Respond(bytesSlice)
			return
		}

		m.Type = pb.WeddingResponse_Success
		m.Wedding = wedding.ToProto()

		logging.WithContext(ctx).With("weddingID", wedding.ID).Info("responding with wedding")

		bytesSlice, _ := proto.Marshal(m)

		msg.Respond(bytesSlice)
	})
}

func respondGob(ctx context.Context, msg *nats.Msg, m interface{}) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(m); err != nil {
		logging.WithContext(ctx).Error(err)
	}

	msg.Respond(buf.Bytes())
}
