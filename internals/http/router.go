package http

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// router links each route to some handler
func router(db *sqlx.DB) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	// This authentication searches if the user is valid in the Toggl API
	e.Use(TogglAuthenticationRequired())
	handler := &Handler{db}

	e.GET("integrations/api/v1/slack", handler.ListIntegrations)
	return e
}

func ServerRunner(db *sqlx.DB, wg sync.WaitGroup) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router(db),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	wg.Done()
}