package http

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/lucastamoios/integrations/internals/storage"
)

// router links each route to some handler
func router(db *sqlx.DB, cache storage.HashStorage) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	// This authentication searches if the user is valid in the Toggl API
	e.Use(TogglAuthenticationRequired(cache))
	handler := &Handler{cache, db}

	e.GET("integrations/api/v1/slack", handler.ListIntegrations)
	return e
}

func ServerRunner(db *sqlx.DB, cache storage.HashStorage, wg sync.WaitGroup) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router(db, cache),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	wg.Done()
}