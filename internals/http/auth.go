package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lucastamoios/integrations/internals/storage"
	"github.com/lucastamoios/integrations/internals/toggl"
)

// TogglAuthenticationRequired is a middleware for authentication that uses the Toggl API to
// check if the user is valid or not
func TogglAuthenticationRequired(cache storage.HashStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Basic ")
		// If the token is already in the cache it is not necessary to check again
		if _, ok := cache.Get(token); ok {
			c.Next()
			return
		}

		client := toggl.New(token)
		_, err := client.GetUser()
		if err == toggl.ErrorUnauthorized || err == toggl.ErrorForbidden {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user token is not valid in Toggl API"})
			c.Abort()
			return
		}

		cache.Set(token, "ok")  // TODO: set with expiration
		c.Set("toggl_token", token)
		c.Next()
	}
}
