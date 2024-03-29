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
	IntegrationID      int64      `json:"integration_id" db:"integration_id"`
	TogglCredentials   string     `json:"toggl_credentials" db:"toggl_credentials"`
	ServiceCredentials string     `json:"service_credentials" db:"service_credentials"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Rules struct {
	ProjectName  string `db:"project"`
	Emoji        string `db:"emoji"`
	DoNotDisturb bool   `db:"do_not_disturb"`
}

func IntegrationRunner(db storage.Database, wg sync.WaitGroup) {
	log.Println("starting Integration Runner")
	for {
		var integrations []Integration
		err := db.Select(&integrations, "get-all-integrations")
		if err != nil {
			log.Fatal("Database error: ", err)
			break
		}
		for _, integration := range integrations {
			var slackRules []Rules
			err = db.Select(&slackRules, "get-rules", integration.IntegrationID)
			if err != nil {
				log.Fatal("Database error: ", err)
				break
			}

			err = updateStatus(integration, slackRules)
			if err != nil && !errors.Is(err, toggl.TogglAPIError{}) {
				log.Println(fmt.Errorf("failed to set slack status for integratin %d with error: %s", integration.IntegrationID, err.Error()))
			}
		}
		time.Sleep(5 * time.Minute)
	}
	wg.Done()
	log.Println("leaving Integration Runner")
}

func getEmojiForProject(rules []Rules, projectName string) string {
	for _, rule := range rules {
		if strings.ToLower(rule.ProjectName) == strings.ToLower(projectName) {
			return rule.Emoji
		}
	}
	return ""
}

func updateStatus(integration Integration, rules []Rules) error {
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
