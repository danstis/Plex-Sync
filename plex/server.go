package plex

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Host defines the data to be stored for server objects
type Host struct {
	Name      string
	Hostname  string
	Port      int
	Ssl       bool
	Token     string
	TvSection int
}

// CreateURI assembles the URI for an API request
func CreateURI(server Host, path string) string {
	if server.Ssl {
		return fmt.Sprintf("https://%v:%v/%v", server.Hostname, server.Port, path)
	}
	return fmt.Sprintf("http://%v:%v/%v", server.Hostname, server.Port, path)
}

// SearchShow returns all episodes for a given TV Show
func SearchShow(server Host, title string) (Show, error) {
	//TODO: Update this to handle types
	uri := CreateURI(server, fmt.Sprintf("search?type=2&query=%v", title))
	// log.Printf("Performing REST request to %q", uri)
	resp, err := apiRequest("GET", uri, server.Token, nil)
	if err != nil {
		return Show{}, fmt.Errorf("error getting episodes for show %q from server %q", title, server.Name)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Show{}, fmt.Errorf("unexpected HTTP Response %q", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Show{}, fmt.Errorf("error reading response %v", err)
	}

	results := SR{}

	err = xml.Unmarshal(body, &results)
	if err != nil {
		return Show{}, fmt.Errorf("error parsing xml response: %v", err)
	}
	for _, x := range results.Directories {
		if x.Name == title {
			return x, nil
		}
	}
	return Show{}, fmt.Errorf("no show found matching name %q", title)
}

//SR contains search results
type SR struct {
	XMLName     xml.Name `xml:"MediaContainer"`
	Directories []Show   `xml:"Directory"`
}

//Show defines the structure of a Plex TV Show
type Show struct {
	ID           int    `xml:"ratingKey,attr"`
	Name         string `xml:"title,attr"`
	EpisodeCount int    `xml:"leafCount,attr"`
}

//SyncWatchedTv synchronises the watched TV Shows
func SyncWatchedTv(source, destination Host) {
	log.Printf("Syncing watched Tv Shows from %q to %q", source.Name, destination.Name)

	// Return all shows on the source server

	// For each show on the source server, enumerate all episodes

	// If the local show is marked as watched check if the remote episode is watched

	// Scrobble the episode on the remote server

}
