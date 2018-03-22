package plex

import (
	"os"
	"path"
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
