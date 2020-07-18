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

	queries, err := goyesql.ParseFile("db/queries.sql")
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
			var emojiRules []slack.EmojiRule
			err = db.Select(&emojiRules, queries["get-emoji-rules"], integration.IntegrationID)
			if err != nil {
				log.Fatal("db.Exec error: ", err)
			}

			slack.UpdateStatus(integration, emojiRules)
			if err != nil {
				panic(fmt.Errorf("failed to set slack status for integratin %d with error: %s", integration.IntegrationID, err.Error()))
			}
		}
		time.Sleep(1*time.Minute)
		fmt.Println("Updated status")
	}
}