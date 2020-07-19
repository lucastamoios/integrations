package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lucastamoios/integrations/internals/toggl"
	"net/http"
	"strings"
)

func TogglAuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Basic ")

		// TODO we probably want to cache this in Redis or, simpler, in some Redis-like structure
		client := toggl.New(token)
		_, err := client.GetUser()
		if err == toggl.ErrorUnauthorized || err == toggl.ErrorForbidden {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user token is not valid in Toggl API"})
			c.Abort()
			return
		}
		// TODO store the token in the context
		c.Next()
	}
}
