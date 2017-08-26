package plex

import "fmt"

// Host defines the data to be stored for server objects
type Host struct {
	Name     string
	Hostname string
	Port     int
	Ssl      bool
	Token string
}

// CreateURI assembles the URI for an API request
func CreateURI(server Host, path, token string) string {
	if server.Ssl {
		return fmt.Sprintf("https://%v:%v/%v?X-Plex-Token=%v", server.Hostname, server.Port, path, token)
	}
	return fmt.Sprintf("http://%v:%v/%v?X-Plex-Token=%v", server.Hostname, server.Port, path, token)
}

// Episodes returns all episodes for a given TV Show
func Episodes(server Host, title, token string) {
	// uri := CreateURI(server, "search?type=2&query={2}", token)

}
