package slack

import (
	"errors"
	"fmt"
	"github.com/lucastamoios/integrations/internals/storage"
	"log"
	"strings"
	"sync"
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

func IntegrationRunner(db storage.Database, wg sync.WaitGroup) {
	for {
		var integrations []Integration
		err := db.Select(&integrations, "get-integrations")
		if err != nil {
			log.Fatal("Database error: ", err)
			break
		}
		for _, integration := range integrations {
			var emojiRules []EmojiRule
			err = db.Select(&emojiRules, "get-emoji-rules", integration.IntegrationID)
			if err != nil {
				log.Fatal("Database error: ", err)
				break
			}

			err = updateStatus(integration, emojiRules)
			if err != nil && !errors.Is(err, toggl.TogglAPIError{}) {
				log.Println(fmt.Errorf("failed to set slack status for integratin %d with error: %s", integration.IntegrationID, err.Error()))
			}
		}
		time.Sleep(5 * time.Minute)
		fmt.Println("Updated status")
	}
	wg.Done()
}

func getEmojiForProject(rules []EmojiRule, projectName string) string {
	for _, rule := range rules {
		if strings.ToLower(rule.ProjectName) == strings.ToLower(projectName) {
			return rule.Emoji
		}
	}
	return ""
}

func updateStatus(integration Integration, rules []EmojiRule) error {
	api := slack.New(integration.ServiceCredentials)
	client := toggl.New(integration.TogglCredentials)

	currentTE, err := client.GetCurrentTimeEntry()
	if errors.Is(err, toggl.ErrorTimeEntryNotFound) {
		err = api.SetUserCustomStatus("", "", 0)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	project, err := client.GetProject(currentTE.WorkspaceID, currentTE.ProjectID)
	if err != nil {
		return err
	}
	emoji := getEmojiForProject(rules, project.Name)

	err = api.SetUserCustomStatus(currentTE.Description, emoji, 0)
	if err != nil {
		return err
	}
	return nil
}
