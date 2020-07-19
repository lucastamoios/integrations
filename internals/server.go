package main

import (
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/lucastamoios/integrations/internals/http"
	"github.com/lucastamoios/integrations/internals/slack"
)

func main() {
	var wg sync.WaitGroup
	db, err := sqlx.Open("postgres", "dbname=toggl_integrations sslmode=disable")
	if err != nil {
		log.Fatal("sqlx.Open error: ", err)
	}
	wg.Add(2) // Integration + server
	go slack.IntegrationRunner(db, wg)
	go http.ServerRunner(db, wg)
	wg.Wait()

}


