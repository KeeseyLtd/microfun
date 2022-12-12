package routes

import (
	"net/http"

	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/v1/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/nats-io/nats.go"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type V1Routes []Route

func Routes(
	router *mux.Router,
	nc *nats.Conn,
) {
	router.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swagger.yaml")
	})

	opts := middleware.RedocOpts{BasePath: "/v1", SpecURL: "swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	setUpWeddingRoutes(router, nc)

	var routes = V1Routes{
		Route{
			"Get API Spec docs",
			"GET",
			"/docs",
			sh.ServeHTTP,
		},
		Route{
			"Login",
			"POST",
			"/auth/login",
			handlers.Login,
		},
		Route{
			"Refresh",
			"POST",
			"/auth/refresh",
			handlers.Refresh,
		},
	}

	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
}
