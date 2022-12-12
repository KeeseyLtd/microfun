package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/KeeseyLtd/microfun/wedding/internal/app"
	"github.com/KeeseyLtd/microfun/wedding/internal/commands"
	"github.com/KeeseyLtd/microfun/wedding/internal/events"
	"github.com/KeeseyLtd/microfun/wedding/internal/logging"
	"github.com/KeeseyLtd/microfun/wedding/internal/queries"

	"github.com/KeeseyLtd/microfun/wedding/internal/config"
	"github.com/KeeseyLtd/microfun/wedding/internal/storage/postgres"
)

func main() {
	ctx := context.Background()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config := config.LoadConfig(".env")
	storage := postgres.NewStorage(*config)

	app := app.Application{
		Queries: app.Queries{
			GetWedding:    queries.NewGetWeddingHandler(storage),
			GetInvitation: queries.NewGetInvitationHandler(storage),
		},
		Commands: app.Commands{
			CreateWedding: commands.NewCreateWeddingHandler(storage),
		},
	}

	nc, err := events.New("nats:4222")
	if err != nil {
		logging.WithContext(ctx).Fatal("could not connect to nats server")
	}

	go func() {
		nc.SubscribeGetWedding(app.Queries.GetWedding)
		nc.SubscribeGetInvitation(app.Queries.GetInvitation)

		nc.SubscribeCreateWedding(app.Commands.CreateWedding)
	}()

	logging.WithContext(ctx).Info("listening on nats")

	<-done

}
