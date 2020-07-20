package http

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucastamoios/integrations/internals/storage"
)

// router links each route to some handler
func router(db storage.Database, cache storage.HashStorage) http.Handler {
	handler := &Handler{cache, db}
	e := gin.New()
	e.Use(gin.Recovery())
	authenticated := e.Group("integrations/api/v1/slack")
	public := e.Group("integrations/api/v1/slack")

	// This authentication searches if the user is valid in the Toggl API
	authenticated.Use(TogglAuthenticationRequired(cache))

	authenticated.GET("/", handler.ListIntegrations)
	authenticated.GET("/setup", handler.SetupSlackIntegration)
	// This is called as callback by external services, so it will not authenticate
	//the user as we don't use any kind of session
	public.GET("/callback", handler.CallbackSetupSlackIntegration)

	return e
}

func ServerRunner(db storage.Database, cache storage.HashStorage, wg sync.WaitGroup) {
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