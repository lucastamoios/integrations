package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres@db/toggl_integrations?sslmode=disable"
	}

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		log.Fatal("migrate.New error: ", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal("m.Up error: ", err)
	}
}
