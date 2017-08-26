package plex

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Host defines the data to be stored for server objects
type Host struct {
	Name     string
	Hostname string
	Port     int
	Ssl      bool
	Token    string
}

// CreateURI assembles the URI for an API request
func CreateURI(server Host, path string) string {
	if server.Ssl {
		return fmt.Sprintf("https://%v:%v/%v", server.Hostname, server.Port, path)
	}
	return fmt.Sprintf("http://%v:%v/%v", server.Hostname, server.Port, path)
}

// Search returns all episodes for a given TV Show
func Search(server Host, title string) (string, error) {
	uri := CreateURI(server, fmt.Sprintf("search?type=2&query=%v", title))
	resp, err := apiRequest("GET", uri, server.Token, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to get episodes for show %q from server %q", title, server.Name)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed with HTTP Response %q", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read body content")
	}

	return string(body), nil
}
