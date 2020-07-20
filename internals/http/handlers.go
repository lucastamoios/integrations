package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"

	"github.com/lucastamoios/integrations/internals/slack"
	"github.com/lucastamoios/integrations/internals/storage"
)

type Handler struct {
	cache storage.HashStorage
	db *sqlx.DB
}

func (h *Handler) ListIntegrations(c *gin.Context) {
	token, ok := c.Get("toggl_token")
	if !ok {
		log.Fatal("Token not found for request")
	}
	queries, err := goyesql.ParseFile("db/queries.sql")
	if err != nil {
		log.Fatal("goyesql.ParseFile error: ", err)
	}
	var integrations []slack.Integration
	err = h.db.Select(&integrations, queries["get-integrations-for-user"], token)
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

func (h *Handler) SetupSlackIntegration(c *gin.Context) {
	temp := make([]byte, 20)
	_, err := rand.Read(temp)
	if err != nil {
		// return err
		log.Fatal("")
	}
	state := base64.URLEncoding.EncodeToString(temp)
	// TODO set expiry
	token := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Basic ")
	h.cache.Set(state, token)

	clientID := os.Getenv("CLIENT_ID")
	scope := "users.profile:write"
	// TODO URL should be parameterized
	redirectURL := "http://localhost:8080/integrations/api/v1/slack/callback"
	url := fmt.Sprintf("https://slack.com/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		clientID,
		scope,
		redirectURL,
		state)
	// TODO redirect user to url
	c.JSON(http.StatusOK, gin.H{"url": url})
}


func (h *Handler) CallbackSetupSlackIntegration(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// If user already have this state we understand as he is the right user
	token, ok := h.cache.Get(state)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saved state is different from what was passed"})
		c.Abort()
		return
	}
	h.cache.Del(state)

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	// TODO URL should be parameterized
	redirectURL := "http://localhost:8080/integrations/api/v1/slack/callback"

	url := fmt.Sprintf("https://slack.com/api/oauth.access?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		clientID,
		clientSecret,
		redirectURL,
		code)
	unpacked, err := makeExternalRequest(url)
	if err != nil {
		log.Printf("Error while making external request for authentication: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "the request for slack server could not be fulfilled"})
		c.Abort()
		return
	}

	queries, err := goyesql.ParseFile("db/queries.sql")
	if err != nil {
		log.Fatal("goyesql.ParseFile error: ", err)
	}
	log.Println("starting query")
	_, err = h.db.Exec(queries["create-integration"], "toggl-slack-integration", token, unpacked["access_token"])
	if err != nil {
		log.Fatal("sql error: ", err)
	}
	log.Println(token)
	// TODO Should redirect also
	c.JSON(http.StatusOK, gin.H{})

}

func makeExternalRequest(url string) (map[string]interface{}, error) {
	var unpacked map[string]interface{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return unpacked, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return unpacked, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return unpacked, err
	}
	err = json.Unmarshal(body, &unpacked)
	if err != nil {
		return unpacked, err
	}
	return unpacked, nil
}
