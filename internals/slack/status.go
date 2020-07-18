package slack

import (
	"strings"
	"time"

	"github.com/lucastamoios/integrations/internals/toggl"

	"github.com/slack-go/slack"
)

type Integration struct {
	IntegrationID int64 `json:"integration_id" db:"integration_id"`
	TogglCredentials string `json:"toggl_credentials" db:"toggl_credentials"`
	ServiceCredentials string `json:"service_credentials" db:"service_credentials"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type EmojiRule struct {
	ProjectName string `db:"project"`
	Emoji string `db:"emoji"`
}

func getEmojiForProject(rules []EmojiRule, projectName string) (string){
	for _, rule := range rules {
		if strings.ToLower(rule.ProjectName) == strings.ToLower(projectName) {
			return rule.Emoji
		}
	}
	return ""
}

func UpdateStatus(integration Integration, rules []EmojiRule) error {
	client := toggl.New(integration.TogglCredentials)
	currentTE, err := client.GetCurrentTimeEntry()
	if err != nil || currentTE == nil {
		return err
	}

	project, err := client.GetProject(currentTE.WorkspaceID, currentTE.ProjectID)
	if err != nil {
		return err
	}
	emoji := getEmojiForProject(rules, project.Name)

	api := slack.New(integration.ServiceCredentials)
	err = api.SetUserCustomStatus(currentTE.Description, emoji, 0)
	if err != nil {
		return err
	}
	return nil
}
