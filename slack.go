package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
	"github.com/slack-go/slack"
)

type Integration struct {
	ID int `json:"id" db:"id"`
	TogglCredentials string `json:"toggl_credentials" db:"toggl_credentials"`
	ServiceCredentials string `json:"service_credentials" db:"service_credentials"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type CurrentTimeEntry struct {
	ID int `json:"id"`
	ProjectID int `json:"project_id"`
	Description string `json:"description"`
}

type SlackProfile struct {
	StatusText string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
	StatusExpiration int `json:"status_expiration"`
}

type ProfilePayload struct {
	Profile SlackProfile `json:"profile"`
}

func getTogglCurrentTimeEntry(token string) (*CurrentTimeEntry, error) {
	url := "https://toggl.com/api/v9/me/time_entries/current"
	auth := fmt.Sprintf("Basic %s", token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", auth)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var currentTE CurrentTimeEntry
	err = json.Unmarshal(body, &currentTE)
	if err != nil {
		return nil, err
	}
	return &currentTE, nil
}

func setSlackStatus(currentTE CurrentTimeEntry, token string) error {
	api := slack.New(token)
	return api.SetUserCustomStatus(currentTE.Description, ":stuck_out_tongue_winking_eye:", 10)
}

func updateUserSlackStatus(togglToken, slackToken string) error {
	currentTE, err := getTogglCurrentTimeEntry(togglToken)
	if err != nil || currentTE == nil {
		return err
	}
	err = setSlackStatus(*currentTE, slackToken)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := sqlx.Open("postgres", "dbname=toggl_integrations sslmode=disable")
	if err != nil {
		log.Fatal("sqlx.Open error: ", err)
	}

	queries, err := goyesql.ParseFile("db/get_integrations.sql")
	if err != nil {
		log.Fatal("goyesql.ParseFile error: ", err)
	}

	for {
		var integrations []Integration
		err = db.Select(&integrations, queries["get-integrations"])
		if err != nil {
			log.Fatal("db.Exec error: ", err)
		}

		for _, integration := range integrations {
			updateUserSlackStatus(integration.TogglCredentials, integration.ServiceCredentials)
			if err != nil {
				panic(fmt.Errorf("failed to set slack status for integratin %d with error: %s", integration.ID, err.Error()))
			}
		}
		time.Sleep(1*time.Minute)
	}
}