package toggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

// This should change according to environment
const APIURL = "https://toggl.com/"

type Client struct {
	apiURL string
	token string
}

type TimeEntry struct {
	ID int `json:"id"`
	ProjectID int `json:"project_id"`
	Description string `json:"description"`
}

func New(token string) *Client {
	return &Client{apiURL: APIURL, token: token}
}

func (c *Client) GetCurrentTimeEntry() (*TimeEntry, error) {
	url := path.Join(c.apiURL, "api/v9/me/time_entries/current")
	auth := fmt.Sprintf("Basic %s", c.token)
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
	var currentTE TimeEntry
	err = json.Unmarshal(body, &currentTE)
	if err != nil {
		return nil, err
	}
	return &currentTE, nil
}