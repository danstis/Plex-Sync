package plex

import "fmt"

// Host defines the data to be stored for server objects
type Host struct {
	Name     string
	Hostname string
	Port     int
	Ssl      bool
}

// CreateURI assembles the URI for an API request
func CreateURI(server Host, path, token string) string {
	if server.Ssl {
		return fmt.Sprintf("https://%v:%v/%v?X-Plex-Token=%v", server.Hostname, server.Port, path, token)
	}
	return fmt.Sprintf("http://%v:%v/%v?X-Plex-Token=%v", server.Hostname, server.Port, path, token)
}
