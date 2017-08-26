package plex

import (
	"testing"
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
		{"Test SSL", args{server: Host{Name: "Test SSL", Hostname: "localhost", Port: 2121, Ssl: true}, path: "test"}, "https://localhost:2121/test"},
		{"Test HTTP", args{server: Host{Name: "Test HTTP", Hostname: "servername", Port: 1515, Ssl: false}, path: "new"}, "http://servername:1515/new"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateURI(tt.args.server, tt.args.path); got != tt.want {
				t.Errorf("CreateURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
