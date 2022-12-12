package tests

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithPostgres() (context.Context, testcontainers.Container) {
	logger := log.New(ioutil.Discard, "prefix", log.LstdFlags)
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgis/postgis",
		ExposedPorts: []string{"5432/tcp", "5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "wedding-qr",
			"POSTGRES_PASSWORD": "wedding-qr",
			"POSTGRES_DB":       "wedding-service",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           logger,
	})

	if err != nil {
		log.Fatal(err)
	}

	p, _ := postgresC.MappedPort(ctx, "5432/tcp")
	port := p.Int()

	viper.Set("DB_PORT", port)

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		"wedding-qr",
		"wedding-qr",
		"wedding-service",
		"localhost",
		port,
		"disable",
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "../../sql/migrations/",
	}
	migrate.SetTable("migrations")

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Unable to migrate database: %v", err)
	}

	file, err := ioutil.ReadFile("../../sql/seed.sql")
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	requests := strings.Split(string(file), ";\n")
	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			log.Fatalf("unable to seed database: %v", err)
		}
	}

	return ctx, postgresC
}
