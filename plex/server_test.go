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
	CacheLifetime = 5

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
		result, err := SearchShow(host, "GoodShow")
		if err != nil {
			t.Errorf("unable to search show: %v", err)
		}
		expected := Show{
			ID:           123,
			Name:         "GoodShow",
			EpisodeCount: 2,
			Thumbnail:    "/library/metadata/123/thumb/123456789",
			Banner:       "/library/metadata/123/banner/123456789",
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

func TestGettingShowEpisodes(t *testing.T) {
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

	Convey("When searching for all show episodes", t, func() {
		result, err := allEpisodes(host, 123)
		if err != nil {
			t.Errorf("unable to get show episodes: %v", err)
		}
		expected := []Episode{
			Episode{ID: 128, Name: "Episode1", Episode: 1, Season: 1, ViewCount: 1, LastWatched: 1519792250},
			Episode{ID: 125, Name: "Episode2", Episode: 2, Season: 1, ViewCount: 0, LastWatched: 0},
		}

		Convey("The correct episode details should be returned", func() {
			So(result, ShouldResemble, expected)
		})

		Convey("And asked to select an existsing episode", func() {
			ep, err := findEpisode(result, 1, 1)
			if err != nil {
				t.Errorf("failed to get single episode from all eisodes: %v", err)
			}
			expectedEp := Episode{ID: 128, Name: "Episode1", Episode: 1, Season: 1, ViewCount: 1, LastWatched: 1519792250}

			Convey("The correct episode is returned", func() {
				So(ep, ShouldResemble, expectedEp)
			})
		})

		Convey("And asked to select a non-existsing episode", func() {
			ep, err := findEpisode(result, 2, 8)
			expectedEp := Episode{}

			Convey("An empty Episode and an error is returned", func() {
				So(ep, ShouldResemble, expectedEp)
				So(err, ShouldResemble, fmt.Errorf("could not find episode on destination server"))
			})
		})
	})

	Convey("When searching for a non-existing show", t, func() {
		result, err := allEpisodes(host, 000)
		expected := []Episode{}

		Convey("An empty show and an error should be returned", func() {
			So(result, ShouldResemble, expected)
			So(err, ShouldResemble, fmt.Errorf("error parsing xml response: %v", "EOF"))
		})
	})
}

func TestMediaScrobbling(t *testing.T) {
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

	Convey("When attempting to scobble an existing episode", t, func() {
		err := scrobble(host, 125)

		Convey("No error should be raised", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("When attempting to scobble a non-existing episode", t, func() {
		err := scrobble(host, 666)

		Convey("An error should be raised", func() {
			So(err, ShouldResemble, fmt.Errorf("unexpected HTTP Response: 500 Internal Server Error"))
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

		case "/search?type=2&query=GoodShow":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <MediaContainer size="1" identifier="com.plexapp.plugins.library" mediaTagPrefix="/system/bundle/media/flags/" mediaTagVersion="1518751453">
                    <Directory allowSync="1" librarySectionID="2" librarySectionTitle="TV Shows" personal="1" ratingKey="123" key="/library/metadata/123/children" studio="A Network" type="show" title="GoodShow" contentRating="TV-MA" summary="Description Blurb." index="1" rating="9.0" viewCount="5" lastViewedAt="1519811198" year="1989" thumb="/library/metadata/123/thumb/123456789" art="/library/metadata/123/art/123456789" banner="/library/metadata/123/banner/123456789" theme="/library/metadata/123/theme/123456789" duration="1500000" originallyAvailableAt="1989-03-01" leafCount="2" viewedLeafCount="2" childCount="1" addedAt="1519355149" updatedAt="123456789">
                        <Genre tag="Reality" />
                    </Directory>
                </MediaContainer>`)

		case "/search?type=2&query=BadShow":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <MediaContainer size="0" identifier="com.plexapp.plugins.library" mediaTagPrefix="/system/bundle/media/flags/" mediaTagVersion="1518751453"></MediaContainer>`)

		case "/library/metadata/123/allLeaves?":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <MediaContainer size="2" allowSync="1" art="/library/metadata/123/art/123456789" banner="/library/metadata/123/banner/123456789" identifier="com.plexapp.plugins.library" key="123" librarySectionID="2" librarySectionTitle="TV Shows" mediaTagPrefix="/system/bundle/media/flags/" mediaTagVersion="1518751453" mixedParents="1" nocache="1" parentIndex="1" parentTitle="GoodShow" parentYear="1989" theme="/library/metadata/123/theme/123456789" title1="TV Shows" title2="GoodShow" viewGroup="episode" viewMode="65592">
                    <Video ratingKey="128" key="/library/metadata/128" parentRatingKey="124" studio="A Network" type="episode" title="Episode1" parentKey="/library/metadata/124" grandparentTitle="GoodShow" parentTitle="Season 1" contentRating="TV-MA" summary="Episode summary." index="1" parentIndex="1" viewCount="1" lastViewedAt="1519792250" year="2018" thumb="/library/metadata/128/thumb/1519862671" art="/library/metadata/-1/art/1519355672" parentThumb="/library/metadata/124/thumb/123456789" grandparentThumb="/library/metadata/-1/thumb/1519355672" grandparentArt="/library/metadata/-1/art/1519355672" grandparentTheme="/library/metadata/-1/theme/1519355672" duration="1194894" originallyAvailableAt="2018-01-22" addedAt="1519355672" updatedAt="1519862671">
                        <Media videoResolution="720" id="1001" duration="1194894" bitrate="4885" width="1280" height="720" aspectRatio="1.78" audioChannels="2" audioCodec="ac3" videoCodec="h264" container="mkv" videoFrameRate="60p" videoProfile="high">
                            <Part id="1001" key="/library/parts/1001/1516684040/file.mkv" duration="1194894" file="/volume1/Media/Series/GoodShow/Season 1/GoodShow - S01E01 - Episode1.mkv" size="729690742" container="mkv" videoProfile="high" />
                        </Media>
                    </Video>
                    <Video ratingKey="125" key="/library/metadata/125" parentRatingKey="124" studio="A Network" type="episode" title="Episode2" parentKey="/library/metadata/124" grandparentTitle="GoodShow" parentTitle="Season 1" contentRating="TV-MA" summary="Episode summary." index="2" parentIndex="1" year="2018" thumb="/library/metadata/125/thumb/1519862676" art="/library/metadata/-1/art/1519355149" parentThumb="/library/metadata/124/thumb/123456789" grandparentThumb="/library/metadata/-1/thumb/1519355149" grandparentArt="/library/metadata/-1/art/1519355149" grandparentTheme="/library/metadata/-1/theme/1519355149" duration="1215335" originallyAvailableAt="2018-01-29" addedAt="1519355149" updatedAt="1519862676">
                        <Media videoResolution="1080" id="1002" duration="1215335" bitrate="3596" width="1920" height="1080" aspectRatio="1.78" audioChannels="2" audioCodec="aac" videoCodec="hevc" container="mkv" videoFrameRate="NTSC" audioProfile="lc" videoProfile="main">
                            <Part id="1002" key="/library/parts/1002/1518352670/file.mkv" duration="1215335" file="/volume1/Media/Series/GoodShow/Season 1/GoodShow - S01E02 - Episode2.mkv" size="546330868" audioProfile="lc" container="mkv" videoProfile="main" />
                        </Media>
                    </Video>
                </MediaContainer>`)

			// 125 is the test for a successful scrobble
		case "/:/scrobble?key=125&identifier=com.plexapp.plugins.library":
			w.WriteHeader(http.StatusOK)

			// 666 is the test for a failed scrobble
		case "/:/scrobble?key=666&identifier=com.plexapp.plugins.library":
			w.WriteHeader(http.StatusInternalServerError)
		}

	}))

	return s
}
