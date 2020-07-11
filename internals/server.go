package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"

	"github.com/lucastamoios/integrations/internals/slack"
)

func main() {
	db, err := sqlx.Open("postgres", "dbname=toggl_integrations sslmode=disable")
	if err != nil {
		log.Fatal("sqlx.Open error: ", err)
	}

	queries, err := goyesql.ParseFile("db/get_integrations.sql")
	if err != nil {
		log.Fatal("goyesql.ParseFile error: ", err)
	}

	for {
		var integrations []slack.Integration
		err = db.Select(&integrations, queries["get-integrations"])
		if err != nil {
			log.Fatal("db.Exec error: ", err)
		}

		for _, integration := range integrations {
			slack.UpdateUserSlackStatus(integration.TogglCredentials, integration.ServiceCredentials)
			if err != nil {
				panic(fmt.Errorf("failed to set slack status for integratin %d with error: %s", integration.ID, err.Error()))
			}
		}
		time.Sleep(1*time.Minute)
	}
}