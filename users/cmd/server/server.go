package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KeeseyLtd/microfun/users/internal/api"
	mid "github.com/KeeseyLtd/microfun/users/internal/api/middleware"
	"github.com/KeeseyLtd/microfun/users/internal/app"
	"github.com/KeeseyLtd/microfun/users/internal/config"
	"github.com/KeeseyLtd/microfun/users/internal/logging"
	"github.com/KeeseyLtd/microfun/users/internal/queries"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Initialize(configLoc string) app.App {
	cfg := config.LoadConfig(configLoc)
	a := app.App{
		Config: cfg,
		Queries: app.Queries{
			GetUserByLogin:        queries.NewGetUserByLoginHandler(cfg),
			GetUserByRefreshToken: queries.NewGetUserByRefreshTokenHandler(cfg),
		},
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:3000"},
		AllowCredentials: true,
	}))
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		mid.SetJsonContentType,
		middleware.Recoverer,
	)

	api.Routes(
		r,
		a,
	)

	a.Router = r

	return a
}

// Run starts the API server
func Run(a app.App) {
	ctx := context.Background()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Handler:      a.Router,
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
