package events_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"testing"
	"time"

	"github.com/KeeseyLtd/microfun/wedding/internal/domain"
	"github.com/KeeseyLtd/microfun/wedding/internal/errs"
	"github.com/KeeseyLtd/microfun/wedding/internal/events"
	"github.com/KeeseyLtd/microfun/wedding/pb"
	"github.com/KeeseyLtd/microfun/wedding/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	nats "github.com/nats-io/nats.go"
)

const TEST_PORT = 8369

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}

const validUUID = "c8822423-e479-4553-ad0c-eb0d46155982"

type weddingSuccessRepoStub struct{}

func (r *weddingSuccessRepoStub) Handle(ctx context.Context, weddingID uuid.UUID) (domain.Wedding, error) {
	return domain.Wedding{
		ID:          uuid.MustParse("c8822423-e479-4553-ad0c-eb0d46155982"),
		WeddingDate: time.Unix(1670773579, 472334500),
	}, nil
}

type weddingNotFoundRepoStub struct{}

func (r *weddingNotFoundRepoStub) Handle(ctx context.Context, weddingID uuid.UUID) (domain.Wedding, error) {
	return domain.Wedding{}, errs.WeddingResourceNotFound
}

func TestNats_SubscribeGetWedding(t *testing.T) {
	type queryer interface {
		Handle(ctx context.Context, weddingID uuid.UUID) (domain.Wedding, error)
	}

	tests := []struct {
		name      string
		n         func() *events.Nats
		want      *pb.Wedding
		arg       string
		handler   queryer
		wantErr   bool
		errorCode errs.ErrorCode
	}{
		{
			name: "OK",
			n: func() *events.Nats {
				nc, _ := events.New(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
				return nc
			},
			arg:     validUUID,
			handler: new(weddingSuccessRepoStub),
			want: &pb.Wedding{
				Uuid:        validUUID,
				WeddingDate: timestamppb.New(time.Unix(1670773579, 472334500)),
			},
		},
		{
			name: "Not Found",
			n: func() *events.Nats {
				nc, _ := events.New(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
				return nc
			},
			arg:       "aa822423-e479-4553-ad0c-eb0d46155982",
			handler:   new(weddingNotFoundRepoStub),
			want:      nil,
			wantErr:   true,
			errorCode: errs.WeddingResourceNotFound,
		},
		{
			name: "Invalid UUID",
			n: func() *events.Nats {
				nc, _ := events.New(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
				return nc
			},
			arg:       "asdadsad",
			handler:   new(weddingSuccessRepoStub),
			want:      nil,
			wantErr:   true,
			errorCode: errs.InvalidUUIDErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := RunServerOnPort(TEST_PORT)
			defer s.Shutdown()

			tt.n().SubscribeGetWedding(tt.handler)

			nc, _ := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
			defer nc.Close()

			id := &pb.GetWedding{
				Uuid: tt.arg,
			}

			byteSlice, _ := proto.Marshal(id)

			msg, err := nc.Request("wedding.get", byteSlice, 5*time.Second)
			if err != nil {
				fmt.Println(err)
			}

			m := pb.WeddingResponse{}
			proto.Unmarshal(msg.Data, &m)

			assert.Equal(t, tt.want, m.Wedding)

			if tt.wantErr {
				assert.Equal(t, pb.WeddingResponse_Error, m.Type)
				assert.Equal(t, int32(tt.errorCode), m.ErrorCode)

				return
			}

			assert.Equal(t, pb.WeddingResponse_Success, m.Type)
		})
	}
}

type invitationSuccessRepoStub struct{}

func (r *invitationSuccessRepoStub) Handle(ctx context.Context, args domain.GetInvitation) (domain.Invitation, error) {
	return domain.Invitation{
		ID: uuid.MustParse("c8822423-e479-4553-ad0c-eb0d46155982"),
	}, nil
}

type invitationNotFoundRepoStub struct{}

func (r *invitationNotFoundRepoStub) Handle(ctx context.Context, args domain.GetInvitation) (domain.Invitation, error) {
	return domain.Invitation{}, errs.InviteResourceNotFound
}

func TestNats_SubscribeGetInvitation(t *testing.T) {
	type queryer interface {
		Handle(ctx context.Context, args domain.GetInvitation) (domain.Invitation, error)
	}

	tests := []struct {
		name      string
		n         func() *events.Nats
		want      domain.Invitation
		arg       domain.GetInvitation
		handler   queryer
		wantErr   bool
		errorCode errs.ErrorCode
	}{
		{
			name: "OK",
			n: func() *events.Nats {
				nc, _ := events.New(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
				return nc
			},
			arg: domain.GetInvitation{
				WeddingID:    uuid.New(),
				InvitationID: uuid.MustParse(validUUID),
			},
			handler: new(invitationSuccessRepoStub),
			want: domain.Invitation{
				ID: uuid.MustParse(validUUID),
			},
		},
		{
			name: "Not Found",
			n: func() *events.Nats {
				nc, _ := events.New(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
				return nc
			},
			arg: domain.GetInvitation{
				WeddingID:    uuid.New(),
				InvitationID: uuid.MustParse("aa822423-e479-4553-ad0c-eb0d46155982"),
			},
			handler:   new(invitationNotFoundRepoStub),
			want:      domain.Invitation{},
			wantErr:   true,
			errorCode: errs.InviteResourceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := RunServerOnPort(TEST_PORT)
			defer s.Shutdown()

			tt.n().SubscribeGetInvitation(tt.handler)

			nc, _ := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
			defer nc.Close()

			var sBuf bytes.Buffer
			gob.NewEncoder(&sBuf).Encode(tt.arg)

			msg, err := nc.Request("wedding.invitation.get", sBuf.Bytes(), 5*time.Second)
			if err != nil {
				fmt.Println(err)
			}

			m := events.GetInvitationResponseMessage{}
			rBuf := bytes.NewBuffer(msg.Data)
			gob.NewDecoder(rBuf).Decode(&m)

			assert.Equal(t, tt.want, m.Payload)

			if tt.wantErr {
				assert.Equal(t, events.Error, m.Type)
				assert.Equal(t, tt.errorCode, m.ErrorCode)

				return
			}

			assert.Equal(t, events.Success, m.Type)
		})
	}
}

func TestNats_SubscribeCreateWedding(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	m := mocks.NewMockcreateWeddingCommand(mockCtrl)

	tests := []struct {
		name    string
		n       func() *events.Nats
		mock    func(domain.CreateWeddingParams)
		arg     domain.CreateWeddingParams
		wantErr bool
	}{
		{
			name: "OK",
			n: func() *events.Nats {
				nc, _ := events.New(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
				return nc
			},
			arg: domain.CreateWeddingParams{
				Names:       "test",
				WeddingDate: time.Now().AddDate(0, -1, 0),
				UserID:      uuid.New(),
			},
			mock: func(args domain.CreateWeddingParams) {
				m.EXPECT().Handle(gomock.Any(), args).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := RunServerOnPort(TEST_PORT)
			defer s.Shutdown()
			tt.mock(tt.arg)

			tt.n().SubscribeCreateWedding(m)

			nc, _ := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT))
			defer nc.Close()

			var sBuf bytes.Buffer
			gob.NewEncoder(&sBuf).Encode(tt.arg)

			nc.Publish("wedding.create", sBuf.Bytes())
			time.Sleep(100 * time.Millisecond)
		})
	}
}
