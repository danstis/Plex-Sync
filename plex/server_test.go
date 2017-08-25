package plex

import (
	"testing"
)

func TestCreateURI(t *testing.T) {
	type args struct {
		server Host
		path   string
		token  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test SSL", args{server: Host{Name: "Test SSL", Hostname: "localhost", Port: 2121, Ssl: true}, path: "test", token: "123456"}, "https://localhost:2121/test?X-Plex-Token=123456"},
		{"Test HTTP", args{server: Host{Name: "Test HTTP", Hostname: "servername", Port: 1515, Ssl: false}, path: "new", token: "789456"}, "http://servername:1515/new?X-Plex-Token=789456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateURI(tt.args.server, tt.args.path, tt.args.token); got != tt.want {
				t.Errorf("CreateURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
