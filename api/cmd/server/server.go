package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/middleware"
	v1 "github.com/KeeseyLtd/microfun/api-gateway/internal/api/v1/routes"
	nats "github.com/nats-io/nats.go"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/config"
	"github.com/KeeseyLtd/microfun/api-gateway/internal/logging"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
	Config *config.Config
}

func (a *App) Initialize(configLoc string) {
	a.Config = config.LoadConfig(configLoc)

	nc, err := nats.Connect(fmt.Sprintf("%s:%s", "nats", "4222"))
	if err != nil {
		logging.WithContext(context.Background()).Error(err)
	}
	// storage := postgres.NewStorage(a.Config.Database)
	// mailer := mailgun.NewMailgunMailer(a.Config.Mail)

	// var cache caching.Cacher
	// switch a.Config.Cache.Type {
	// case "redis":
	// 	cache = redis.New(a.Config.Cache)
	// case "memcache":
	// 	cache = memcache.New(a.Config.Cache)
	// default:
	// 	cache = nil
	// 	a.Config.Cache.Enabled = false
	// }

	// qrService := qrs.NewService(storage, a.Config)
	// weddingService := weddings.NewService(storage, a.Config, mailer, cache)
	// userService := users.NewService(storage, a.Config, mailer)

	r := mux.NewRouter().StrictSlash(true)
	v1.Routes(
		r.PathPrefix("/v1").Subrouter(),
		nc,
		// qrService,
		// weddingService,
		// userService,
	)

	r.Use(
		middleware.AddRequestUUID,
		middleware.LogRequests,
		middleware.SetJsonContentType,
		middleware.AuthCheck,
	)

	a.Router = r
}

// Run starts the API server
func (a *App) Run() {
	ctx := context.Background()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:3000"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(a.Router),
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", a.Config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.WithContext(ctx).With("error", err).Fatal("failed to start server")
		}
	}()

	logging.WithContext(ctx).Infof("started server %s:%d", "0.0.0.0", a.Config.Port)

	<-done

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.WithContext(ctx).With("error", err).Fatal("failed to shutdown server")
	}

	logging.WithContext(ctx).Infof("server exited properly")
}
