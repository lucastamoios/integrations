package http

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucastamoios/integrations/internals/slack"
	"github.com/nleof/goyesql"
)

type Handler struct {
	db *sqlx.DB
}

func (h *Handler) ListIntegrations(c *gin.Context) {
	queries, err := goyesql.ParseFile("db/queries.sql")
	if err != nil {
		log.Fatal("goyesql.ParseFile error: ", err)
	}
	var integrations []slack.Integration
	err = h.db.Select(&integrations, queries["get-integrations"])
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"error": "Some problem happened",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"integrations": integrations,
		},
	)
}
