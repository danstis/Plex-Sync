package plex

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestCreateURI(t *testing.T) {
	type args struct {
		server Host
		path   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test SSL",
			args: args{server: Host{Name: "Test SSL", Hostname: "localhost", Port: 2121, Ssl: true}, path: "test"},
			want: "https://localhost:2121/test",
		},
		{
			name: "Test HTTP",
			args: args{server: Host{Name: "Test HTTP", Hostname: "servername", Port: 1515, Ssl: false}, path: "new"},
			want: "http://servername:1515/new",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateURI(tt.args.server, tt.args.path); got != tt.want {
				t.Errorf("CreateURI() = %v, want %v", got, tt.want)
			}
		})
	}
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
