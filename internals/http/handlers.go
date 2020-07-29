package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucastamoios/integrations/internals/slack"
	"github.com/lucastamoios/integrations/internals/storage"
)

type Handler struct {
	cache storage.HashStorage
	db    storage.Database
}

type StateInformation struct {
	Token    string
	ReturnTo string
}

var SLACK_INTEGRATION_ROUTE = "integrations/api/v1/slack"
var CALLBACK_SUBROUTE = "callback"

func (h *Handler) ListIntegrations(c *gin.Context) {
	token, ok := c.Get("toggl_token")
	if !ok {
		log.Fatal("token not found for request")
	}
	var integrations []slack.Integration
	err := h.db.Select(&integrations, "get-integrations-for-user", token)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"error": "some problem happened loading the integrations",
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
		log.Fatal("failed generating state")
		c.Abort()
		return
	}

	state := base64.URLEncoding.EncodeToString(temp)
	returnTo := c.Query("return_to")
	token := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Basic ")
	h.cache.Set(state, StateInformation{token, returnTo})
	h.cache.Expire(state, 10*time.Minute)

	clientID := os.Getenv("CLIENT_ID")
	scope := "users.profile:write"
	callbackURL, err := url.Parse(os.Getenv("HOST"))
	if err != nil {
		log.Fatal("failed parsing host")
		c.JSON(http.StatusBadGateway, gin.H{"error": "an internal error occurred"})
		c.Abort()
		return
	}
	callbackURL.Path = path.Join(callbackURL.Path, SLACK_INTEGRATION_ROUTE, CALLBACK_SUBROUTE)

	url := fmt.Sprintf("https://slack.com/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		clientID,
		scope,
		callbackURL.String(),
		state)
	c.Redirect(http.StatusFound, url)
}

func (h *Handler) CallbackSetupSlackIntegration(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// If user already have this state we understand as he is the right user
	stateInformation, ok := h.cache.Get(state)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "saved state is different from what was passed or it expired"})
		c.Abort()
		return
	}
	h.cache.Del(state)

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	callbackURL, err := url.Parse(os.Getenv("HOST"))
	if err != nil {
		log.Fatal("failed parsing host")
		c.JSON(http.StatusBadGateway, gin.H{"error": "an internal error occurred"})
		c.Abort()
		return
	}
	callbackURL.Path = path.Join(callbackURL.Path, SLACK_INTEGRATION_ROUTE, CALLBACK_SUBROUTE)
	url := fmt.Sprintf("https://slack.com/api/oauth.access?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		clientID,
		clientSecret,
		callbackURL.String(),
		code)
	unpacked, err := makeExternalRequest(url)
	if err != nil {
		log.Printf("Error while making external request for authentication: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "the request for slack server could not be fulfilled"})
		c.Abort()
		return
	}

	token := stateInformation.(StateInformation).Token
	returnTo := stateInformation.(StateInformation).ReturnTo
	_, err = h.db.Exec("create-integration", "toggl-slack-integration", token, unpacked["access_token"])
	if err != nil {
		log.Fatal("sql error: ", err)
	}
	c.Redirect(http.StatusFound, returnTo)

}

func (h *Handler) ListSlackRules(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		gin.H{
			"error": "not implemented yet",
		},
	)
}

func (h *Handler) CreateSlackRules(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		gin.H{
			"error": "not implemented yet",
		},
	)
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
