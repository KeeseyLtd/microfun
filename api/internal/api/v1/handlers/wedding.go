package handlers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/apierrors"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/domain"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/logging"

	"github.com/KeeseyLtd/microfun/api-gateway/pb"
)

const (
	GetWeddingTopic    = "wedding.get"
	GetInvitationTopic = "wedding.invitation.get"

	CreateWeddingTopic = "wedding.create"
)

func GetWedding(nc *nats.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["weddingID"]
		parsedID, err := uuid.Parse(id)
		if err != nil {
			logging.WithContext(r.Context()).Infow(apierrors.InvalidUUIDErr.Error(), "error", err)

			ErrorResponse(w, apierrors.APIResponse(apierrors.InvalidUUIDErr))

			return
		}

		wID := &pb.GetWedding{
			Uuid: parsedID.String(),
		}

		byteSlice, _ := proto.Marshal(wID)
		msg, err := nc.Request(GetWeddingTopic, byteSlice, 5*time.Second)
		if err != nil {
			logging.WithContext(r.Context()).With("error", err).Error("did not recieve reply correctly")
			ErrorResponse(w, apierrors.APIResponse(apierrors.SomethingWentWrong))

			return
		}

		m := pb.WeddingResponse{}
		if err := proto.Unmarshal(msg.Data, &m); err != nil {
			logging.WithContext(r.Context()).With("error", err).Error("could not decode message")
			ErrorResponse(w, apierrors.APIResponse(apierrors.SomethingWentWrong))

			return
		}

		if m.Type != pb.WeddingResponse_Success {
			ErrorResponse(w, apierrors.APIResponse(apierrors.ErrorCode(m.GetErrorCode())))

			return
		}

		res := WeddingResponse{}
		res.Body.Data = domain.NewWeddingFromProto(m.Wedding)

		successResponse(w, http.StatusOK, res)
	}
}

func CreateWedding(nc *nats.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var createWeddingForm domain.CreateWeddingForm

		if err := json.NewDecoder(r.Body).Decode(&createWeddingForm); err != nil {
			logging.WithContext(r.Context()).With("error", err).Info(apierrors.InvalidJson)

			ErrorResponse(w, apierrors.APIResponse(apierrors.InvalidJson))
			return
		}

		if !createWeddingForm.Validate() {
			e := apierrors.APIResponse(apierrors.ValidationErr)
			ErrorResponse(w, e.Details(createWeddingForm.Errors))

			return
		}

		weddingDate, _ := time.Parse("2006-01-02 15:04:05", createWeddingForm.WeddingDate)

		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(domain.CreateWeddingParams{
			Names:       createWeddingForm.Names,
			WeddingDate: weddingDate,
			UserID:      r.Context().Value(domain.UserInfoKey).(uuid.UUID),
		}); err != nil {
			logging.WithContext(r.Context()).With("error", err).Error("could not encode wedding params")
			ErrorResponse(w, apierrors.SomethingWentWrong)

			return
		}

		if err := nc.Publish(CreateWeddingTopic, buf.Bytes()); err != nil {
			logging.WithContext(r.Context()).With("error", err).Error("could not publish create wedding event")

			ErrorResponse(w, apierrors.SomethingWentWrong)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
