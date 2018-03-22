package plex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURLCreation(t *testing.T) {

	Convey("When asked to generate an HTTP URL", t, func() {
		uri := CreateURI(Host{Name: "Test SSL", Hostname: "localhost", Port: 2121, Ssl: true}, "test")

		Convey("An HTTP URL should be generated", func() {
			So(uri, ShouldEqual, "https://localhost:2121/test")
		})
	})

	Convey("When asked to generate an HTTPS URL", t, func() {
		uri := CreateURI(Host{Name: "Test HTTP", Hostname: "servername", Port: 1515, Ssl: false}, "new")

		Convey("An HTTPS URL should be generated", func() {
			So(uri, ShouldEqual, "http://servername:1515/new")
		})
	})
}

func TestCachedFileExpiry(t *testing.T) {
	CacheLifetime = -5

	Convey("Given a file created yesterday", t, func() {
		newFile := tempfile(t, "newfile", -1)

		Convey("The file should not be marked as expired", func() {
			So(expired(newFile), ShouldEqual, false)
		})
	})

	Convey("Given a file created 99 days ago", t, func() {
		oldFile := tempfile(t, "oldfile", -99)

		Convey("The file should be marked as expired", func() {
			So(expired(oldFile), ShouldEqual, true)
		})
	})
}

func TestSearchShow(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()
	port, err := strconv.Atoi(strings.Split(ts.URL, ":")[2])
	if err != nil {
		t.Errorf("failed to start test server")
	}
	host := Host{
		Name:     "TestServer",
		Hostname: "127.0.0.1",
		Port:     port,
		Ssl:      false,
	}

	Convey("When searching for as existing show", t, func() {
		result, err := SearchShow(host, "Cops")
		if err != nil {
			t.Errorf("unable to search show: %v", err)
		}
		expected := Show{
			ID:           273,
			Name:         "Cops",
			EpisodeCount: 5,
			Thumbnail:    "/library/metadata/273/thumb/1519355692",
			Banner:       "/library/metadata/273/banner/1519355692",
		}

		Convey("The correct show details should be returned", func() {
			So(result, ShouldResemble, expected)
		})
	})

	Convey("When searching for a non existing show", t, func() {
		result, err := SearchShow(host, "BadShow")
		expected := Show{}

		Convey("An empty show and an error should be returned", func() {
			So(result, ShouldResemble, expected)
			So(err, ShouldResemble, fmt.Errorf("no show found matching name %q", "BadShow"))
		})
	})
}

//Helper functions
func fileInfo(t *testing.T, fn string) os.FileInfo {
	fi, err := os.Stat(fn)
	if err != nil {
		t.Fatal(err)
	}
	return fi
}

func tempfile(t *testing.T, name string, age int) os.FileInfo {
	filepath := path.Join(os.TempDir(), name)
	f, err := os.Create(filepath)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("Test")
	f.Close()
	timestamp := time.Now().AddDate(0, 0, age)
	os.Chtimes(filepath, timestamp, timestamp)
	defer os.Remove(filepath)
	return fileInfo(t, filepath)
}

func startTestServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() + "?" + r.URL.RawQuery {
		case "/search?type=2&query=Cops":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <MediaContainer size="1" identifier="com.plexapp.plugins.library" mediaTagPrefix="/system/bundle/media/flags/" mediaTagVersion="1518751453">
                    <Directory allowSync="1" librarySectionID="2" librarySectionTitle="TV Shows" personal="1" ratingKey="273" key="/library/metadata/273/children" studio="Paramount Network" type="show" title="Cops" contentRating="TV-MA" summary="Description Blurb." index="1" rating="9.0" viewCount="5" lastViewedAt="1519811198" year="1989" thumb="/library/metadata/273/thumb/1519355692" art="/library/metadata/273/art/1519355692" banner="/library/metadata/273/banner/1519355692" theme="/library/metadata/273/theme/1519355692" duration="1500000" originallyAvailableAt="1989-03-01" leafCount="5" viewedLeafCount="5" childCount="1" addedAt="1519355149" updatedAt="1519355692">
                        <Genre tag="Reality" />
                    </Directory>
                </MediaContainer>`)
		case "/search?type=2&query=BadShow":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <MediaContainer size="0" identifier="com.plexapp.plugins.library" mediaTagPrefix="/system/bundle/media/flags/" mediaTagVersion="1518751453"></MediaContainer>`)
		}
	}))

	return s
}
