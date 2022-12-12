package app

import (
	"github.com/KeeseyLtd/microfun/users/internal/config"
	"github.com/KeeseyLtd/microfun/users/internal/queries"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Router  *chi.Mux
	Config  *config.Config
	Queries Queries
}

type Queries struct {
	GetUserByLogin        queries.GetUserByLogin
	GetUserByRefreshToken queries.GetUserByRefreshToken
}
