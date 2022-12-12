package queries_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/KeeseyLtd/microfun/wedding/internal/config"
	"github.com/KeeseyLtd/microfun/wedding/internal/domain"
	"github.com/KeeseyLtd/microfun/wedding/internal/errs"
	"github.com/KeeseyLtd/microfun/wedding/internal/queries"
	"github.com/KeeseyLtd/microfun/wedding/internal/storage/postgres"
	"github.com/KeeseyLtd/microfun/wedding/tests"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	ctx, postgresC := tests.TestWithPostgres()

	m.Run()

	postgresC.Terminate(ctx)
}

func TestGetWedding_Handle(t *testing.T) {
	c := config.LoadConfig("../../.testing.env")
	r := postgres.NewStorage(*c)

	handler := queries.NewGetWeddingHandler(r)

	type args struct {
		ctx       context.Context
		weddingID uuid.UUID
	}
	tests := []struct {
		name      string
		h         queries.GetWedding
		args      args
		want      domain.Wedding
		wantErr   bool
		errorCode errs.ErrorCode
	}{
		{
			name: "OK",
			h:    handler,
			args: args{
				ctx:       context.Background(),
				weddingID: uuid.MustParse("6fd7e79a-45b2-44a4-9308-26e762d63f24"),
			},
			want: domain.Wedding{
				ID:          uuid.MustParse("6fd7e79a-45b2-44a4-9308-26e762d63f24"),
				Names:       "test",
				WeddingDate: timeMustParse("2006-01-02 15:04:05.000").UTC(),
				Invitations: []domain.Invitation{},
			},
		},
		{
			name: "Not Found",
			h:    handler,
			args: args{
				ctx:       context.Background(),
				weddingID: uuid.MustParse("6fd7e79a-aaaa-44a4-9308-26e762d63f24"),
			},
			want:      domain.Wedding{},
			wantErr:   true,
			errorCode: errs.WeddingResourceNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Handle(tt.args.ctx, tt.args.weddingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWedding.Handle() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetWedding.Handle() errorCode = %v, want errorCode %v", err.(errs.ErrorCode), tt.errorCode)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWedding.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func timeMustParse(timeString string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05.000", timeString)
	return t
}
