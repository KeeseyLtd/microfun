// Package classification Wedding QR.
//
// Wedding qr api.
//
//     Schemes: http, https
//     BasePath: /v1
//     Version: 0.0.1
//     Host: localhost:8081
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"github.com/KeeseyLtd/microfun/users/cmd/server"
)

func main() {
	a := server.Initialize(".env")
	server.Run(a)
}
