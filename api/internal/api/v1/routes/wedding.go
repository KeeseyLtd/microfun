package routes

import (
	"github.com/KeeseyLtd/microfun/api-gateway/internal/api/v1/handlers"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func setUpWeddingRoutes(router *mux.Router, nc *nats.Conn) {
	routes := V1Routes{
		// swagger:route GET /weddings/{WeddingIDParam} Weddings GetWedding
		//
		// Gets Wedding details
		//
		//	Produces:
		//		- application/json
		//
		//	Responses:
		//		200: WeddingResponse
		//		401: MessageErr
		Route{
			"Get Wedding",
			"GET",
			"/{weddingID}",
			handlers.GetWedding(nc),
		},

		// swagger:route POST /weddings Weddings CreateWedding

		// Create new wedding

		// 	Produces:
		// 		- application/json

		// 	Responses:
		// 		204: description:No Content
		// 		401: MessageErr
		// 		422: ValidationErr
		Route{
			"Create wedding",
			"POST",
			"",
			handlers.CreateWedding(nc),
		},
	}

	weddingRouter := router.PathPrefix("/weddings").Subrouter()

	for _, route := range routes {
		weddingRouter.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
}
