package plex

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// CacheLifetime controls when a cached image will be refreshed in days
var CacheLifetime int

var cachePath = filepath.Join(".cache", "show")

// Host defines the data to be stored for server objects
type Host struct {
	gorm.Model
	Name     string
	Hostname string
	Port     int
	Ssl      bool
	Token    string
}

var (
	tvShowFile = filepath.Join("config", "tvshows.txt")
)

// CreateURI assembles the URI for an API request
func CreateURI(server Host, path string) string {
	if server.Ssl {
		return fmt.Sprintf("https://%v:%v/%v", server.Hostname, server.Port, path)
	}
	return fmt.Sprintf("http://%v:%v/%v", server.Hostname, server.Port, path)
}

// SearchShow returns all episodes for a given TV Show
func SearchShow(server Host, title string) (Show, error) {
	uri := CreateURI(server, fmt.Sprintf("search?type=2&query=%v", url.PathEscape(title)))
	// log.Printf("Performing REST request to %q", uri)
	resp, err := apiRequest("GET", uri, server.Token, nil)
	if err != nil {
		return Show{}, fmt.Errorf("error getting episodes for show %q from server %q: %v", title, server.Name, err)
	}
	defer resp.Body.Close() // nolint: errcheck

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

// SR contains search results
type SR struct {
	XMLName     xml.Name `xml:"MediaContainer"`
	Directories []Show   `xml:"Directory"`
}

// Show defines the structure of a Plex TV Show
type Show struct {
	ID           int    `xml:"ratingKey,attr"`
	Name         string `xml:"title,attr"`
	EpisodeCount int    `xml:"leafCount,attr"`
	Thumbnail    string `xml:"thumb,attr"`
	Banner       string `xml:"banner,attr"`
}

// ER contains episode results
type ER struct {
	XMLName xml.Name  `xml:"MediaContainer"`
	Video   []Episode `xml:"Video"`
}

// Episode defines the structure of a Plex TV Episode
type Episode struct {
	ID          int    `xml:"ratingKey,attr"`
	Name        string `xml:"title,attr"`
	Episode     int    `xml:"index,attr"`
	Season      int    `xml:"parentIndex,attr"`
	ViewCount   int    `xml:"viewCount,attr"`
	LastWatched int    `xml:"lastViewedAt,attr"`
}

// SyncWatchedTv synchronises the watched TV Shows
func SyncWatchedTv(source, destination Host) error {
	log.Printf("Syncing watched Tv Shows from %q to %q", source.Name, destination.Name)

	// Return all selected shows
	ss, err := SelectedShows()
	if err != nil {
		return err
	}

	// For each show, enumerate all source and destination episodes
	for _, s := range ss {
		log.Printf("Processing show %q", s)
		destShow, err := SearchShow(destination, s)
		if err != nil {
			log.Println(err)
			continue
		}
		err = destShow.cacheImages(destination)
		if err != nil {
			log.Println(err)
		}
		dEps, err := allEpisodes(destination, destShow.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		srcShow, err := SearchShow(source, s)
		if err != nil {
			log.Println(err)
			continue
		}
		sEps, err := allEpisodes(source, srcShow.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, e := range sEps {
			// If the local show is marked as watched check if the remote episode is watched
			log.Printf("- Checking %v - Season %v, Episode %v", srcShow.Name, e.Season, e.Episode)
			destEp, err := findEpisode(dEps, e.Season, e.Episode)
			if err != nil {
				log.Println(err)
				continue
			}
			if e.ViewCount > 0 && destEp.ViewCount < 1 {
				// Scrobble the episode on the remote server
				err = scrobble(destination, destEp.ID)
				if err != nil {
					log.Printf("failed to scrobble episode. Error: %v", err)
					continue
				}
				log.Printf("* Scrobbled on %q", destination.Name)
			} else if destEp.ViewCount >= 1 {
				log.Println("- Already scrobbled, skipping...")
			} else {
				log.Println("- Episode not yet watched, skipping...")
			}
		}
	}
	return nil
}

// SelectedShows returns the selected tv shows from the tvShowsFile
func SelectedShows() ([]string, error) {
	file, err := os.Open(tvShowFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open tvshows file %q", tvShowFile)
	}

	defer file.Close() // nolint: errcheck
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	sort.Strings(lines)
	return lines, scanner.Err()
}

// allEpisodes returns all child episodes of a tv show regardless of the season they belong to
func allEpisodes(server Host, sID int) ([]Episode, error) {
	uri := CreateURI(server, fmt.Sprintf("library/metadata/%v/allLeaves", sID))
	resp, err := apiRequest("GET", uri, server.Token, nil)
	if err != nil {
		return []Episode{}, err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != http.StatusOK {
		return []Episode{}, fmt.Errorf("unexpected HTTP Response %q", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Episode{}, fmt.Errorf("error reading response %v", err)
	}

	results := ER{}

	err = xml.Unmarshal(body, &results)
	if err != nil {
		return []Episode{}, fmt.Errorf("error parsing xml response: %v", err)
	}

	return results.Video, nil
}

// findEpisode returns a single episode from a slice of Episodes based on the season and episode number
func findEpisode(eps []Episode, s, e int) (Episode, error) {
	for _, i := range eps {
		if i.Season == s && i.Episode == e {
			return i, nil
		}
	}

	return Episode{}, fmt.Errorf("could not find episode on destination server")
}

// scrobble marks an episode as watched
func scrobble(server Host, eID int) error {
	uri := CreateURI(server, fmt.Sprintf(":/scrobble?key=%v&identifier=com.plexapp.plugins.library", eID))
	resp, err := apiRequest("GET", uri, server.Token, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP Response: %v", resp.Status)
	}

	return nil
}

// cacheImage downloads an image from the specified server to the cache location
func (s Show) cacheImages(server Host) error {
	itemname := fmt.Sprintf("%s_thumb.jpg", s.Name)
	fullpath := filepath.Join(cachePath, itemname)

	// Check if file is already cached
	if fs, err := os.Stat(fullpath); !os.IsNotExist(err) {
		if !expired(fs) {
			return nil
		}
		log.Println("Cached image is expired, will refresh")
	}

	uri := CreateURI(server, strings.TrimPrefix(s.Thumbnail, "/"))
	resp, err := apiRequest("GET", uri, server.Token, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response %v", err)
	}

	err = os.MkdirAll(cachePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fullpath, body, 0777)
	if err != nil {
		return err
	}
	log.Printf("- Cached banner image to path %q", fullpath)
	return nil
}

func expired(fs os.FileInfo) bool {
	return fs.ModTime().Before(time.Now().AddDate(0, 0, CacheLifetime*-1))
}
