package toggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// This should change according to environment
const APIURL = "https://toggl.com/"

type TogglAPIError struct {
	status int
	message string
	endpoint string
}

func (t TogglAPIError) Error() string {
	return fmt.Sprintf("endpoint %d returned a status %d: %s", t.endpoint, t.status, t.message)
}

var ErrorTimeEntryNotFound = TogglAPIError{404, "not found", APIURL}
var ErrorUnauthorized = TogglAPIError{401, "unauthorized", APIURL}
var ErrorForbidden = TogglAPIError{403, "forbidden", APIURL}

type Client struct {
	apiURL url.URL
	token string
}

type TimeEntry struct {
	ID int `json:"id"`
	WorkspaceID int64 `json:"workspace_id"`
	ProjectID int64 `json:"project_id"`
	Description string `json:"description"`
}

type Project struct {
	ProjectID int64 `json:"id"`
	Name string `json:"name"`
}

type TogglUser struct {
	UserID int64 `json:"id"`
	Name string `json:"fullname"`
	Mail string `json:"email"`
}

func New(token string) *Client {
	u, err := url.Parse(APIURL)
	if err != nil {
		panic(fmt.Errorf("provided URL %s can't be parsed: %s", APIURL, err.Error()))
	}
	return &Client{apiURL: *u, token: token}
}

func (c *Client) GetCurrentTimeEntry() (*TimeEntry, error) {
	c.apiURL.Path = "api/v9/me/time_entries/current"
	body, err := c.makeRequest()
	if err != nil {
		return nil, err
	}

	var currentTE TimeEntry
	err = json.Unmarshal(body, &currentTE)
	if err != nil {
		return nil, err
	}
	return &currentTE, nil
}

func (c *Client) GetProject(workspaceID, projectID int64) (*Project, error) {
	c.apiURL.Path = fmt.Sprintf("api/v9/workspaces/%d/projects/%d", workspaceID, projectID)
	body, err := c.makeRequest()
	if err != nil {
		return nil, err
	}

	var project Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) GetUser() (*TogglUser, error) {
	c.apiURL.Path = "api/v9/me"
	body, err := c.makeRequest()
	if err != nil {
		return nil, err
	}

	var user TogglUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) makeRequest() ([]byte, error) {
	auth := fmt.Sprintf("Basic %s", c.token)
	req, err := http.NewRequest("GET", c.apiURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", auth)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrorUnauthorized
	}
	if res.StatusCode == http.StatusForbidden {
		return nil, ErrorForbidden
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// Some Toggl endpoints return this string instead of a 404 status
	if string(body) == "null" {
		return nil, ErrorTimeEntryNotFound
	}
	return body, nil
}
