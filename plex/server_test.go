package plex

import (
	"reflect"
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

func TestSearchShow(t *testing.T) {
	type args struct {
		server Host
		title  string
	}
	tests := []struct {
		name    string
		args    args
		want    Show
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchShow(tt.args.server, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchShow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchShow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncWatchedTv(t *testing.T) {
	type args struct {
		source      Host
		destination Host
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncWatchedTv(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("SyncWatchedTv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_selectedShows(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectedShows()
			if (err != nil) != tt.wantErr {
				t.Errorf("selectedShows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectedShows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_allEpisodes(t *testing.T) {
	type args struct {
		server Host
		sID    int
	}
	tests := []struct {
		name    string
		args    args
		want    []Episode
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allEpisodes(tt.args.server, tt.args.sID)
			if (err != nil) != tt.wantErr {
				t.Errorf("allEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allEpisodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findEpisode(t *testing.T) {
	type args struct {
		eps []Episode
		s   int
		e   int
	}
	tests := []struct {
		name    string
		args    args
		want    Episode
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findEpisode(tt.args.eps, tt.args.s, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("findEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scrobble(t *testing.T) {
	type args struct {
		server Host
		eID    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := scrobble(tt.args.server, tt.args.eID); (err != nil) != tt.wantErr {
				t.Errorf("scrobble() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
