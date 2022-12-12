package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/KeeseyLtd/microfun/wedding/internal/config"
	"github.com/KeeseyLtd/microfun/wedding/internal/logging"
	_ "github.com/lib/pq"
)

var db *sql.DB

func NewStorage(c config.Config) *Queries {
	var err error
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		c.DBUsername,
		c.DBPassword,
		c.DBDatabase,
		c.DBHostname,
		c.DBPort,
		c.DBSSL,
	)

	s, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := s.Ping(); err != nil {
		logging.WithContext(context.TODO()).Fatalw("Unable to ping the database", "error", err, "host", c.DBHostname)
	}

	db = s
	return New(s)
}

func StartTx() (*sql.Tx, error) {
	return db.Begin()
}
