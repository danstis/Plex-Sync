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

func Test_expired(t *testing.T) {
	type args struct {
		fs os.FileInfo
	}

	// Create some test files
	files := []struct {
		name string
		age  int
	}{
		{"newfile", -1},
		{"oldfile", -50},
	}
	for _, i := range files {
		filepath := path.Join(os.TempDir(), i.name)
		f, err := os.Create(filepath)
		if err != nil {
			t.Fatal("unable to create test file")
		}
		f.WriteString("Test")
		f.Close()
		timestamp := time.Now().AddDate(0, 0, i.age)
		os.Chtimes(filepath, timestamp, timestamp)
		defer os.Remove(filepath)
	}
	CacheLifetime = -5

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test nonexpired file",
			args: args{fileInfo(t, path.Join(os.TempDir(), "newfile"))},
			want: false,
		},
		{
			name: "Test expired file",
			args: args{fileInfo(t, path.Join(os.TempDir(), "oldfile"))},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expired(tt.args.fs); got != tt.want {
				t.Errorf("expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func fileInfo(t *testing.T, fn string) os.FileInfo {
	fi, err := os.Stat(fn)
	if err != nil {
		t.Fatal(err)
	}
	return fi
}
