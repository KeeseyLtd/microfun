package app

import (
	"github.com/KeeseyLtd/microfun/wedding/internal/commands"
	"github.com/KeeseyLtd/microfun/wedding/internal/queries"
)

type Application struct {
	Queries  Queries
	Commands Commands
}

type Queries struct {
	GetWedding    queries.GetWedding
	GetInvitation queries.GetInvitation
}

type Commands struct {
	CreateWedding commands.CreateWedding
}
