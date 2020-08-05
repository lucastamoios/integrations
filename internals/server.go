package main

import (
	"log"
	"sync"

	"github.com/lucastamoios/integrations/internals/http"
	"github.com/lucastamoios/integrations/internals/slack"
	"github.com/lucastamoios/integrations/internals/storage"
)

func main() {
	log.Println("starting Integration Service...")
	var wg sync.WaitGroup
	db, err := storage.NewPostgresDatabase("db/queries.sql")
	if err != nil {
		log.Fatal("NewPostgresDatabase error: ", err)
	}
	cache := storage.NewHashStorage()
	wg.Add(2) // Integration + server
	go slack.IntegrationRunner(db, wg)
	go http.ServerRunner(db, cache, wg)
	log.Println("finished setting service up.")
	wg.Wait()
	log.Println("stopping service.")
}
