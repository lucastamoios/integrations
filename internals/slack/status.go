package slack

import (
	"time"

	"github.com/lucastamoios/integrations/internals/toggl"

	"github.com/slack-go/slack"
)

type Integration struct {
	ID int `json:"id" db:"id"`
	TogglCredentials string `json:"toggl_credentials" db:"toggl_credentials"`
	ServiceCredentials string `json:"service_credentials" db:"service_credentials"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func UpdateUserSlackStatus(togglToken, slackToken string) error {
	client := toggl.New(togglToken)
	currentTE, err := client.GetCurrentTimeEntry()
	if err != nil || currentTE == nil {
		return err
	}
	api := slack.New(slackToken)
	err = api.SetUserCustomStatus(currentTE.Description, ":scream:", 10)
	if err != nil {
		return err
	}
	return nil
}
