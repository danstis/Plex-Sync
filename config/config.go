package config

import (
	"bufio"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/danstis/Plex-Sync/plex"
)

//ConfigFile defines the path to the configuration file
var (
	configFile = path.Join("config", "config.toml")
)

// Settings defines the program configureation
type Settings struct {
	WebServerPort int
	CacheLifetime int
	SyncInterval  int
	Logging       logging   `toml:"logging"`
	LocalServer   plex.Host `toml:"localserver"`
	RemoteServer  plex.Host `toml:"remoteserver"`
}

type logging struct {
	Logfile          string
	Webserverlogfile string
	MaxLogSize       int
	MaxLogCount      int
	MaxLogAge        int
}

// GetConfig returns the application configuration from the config TOML file.
func GetConfig() (Settings, error) {
	var s Settings
	_, err := toml.DecodeFile(configFile, &s)
	if os.IsNotExist(err) {
		s := Settings{
			WebServerPort: 8080,
			Logging: logging{
				MaxLogSize:  5,
				MaxLogCount: 1,
				MaxLogAge:   30,
			},
			LocalServer:   plex.Host{},
			RemoteServer:  plex.Host{},
			CacheLifetime: 30,
			SyncInterval:  3600,
		}
		err := UpdateConfig(s)
		return s, err
	}
	if err != nil {
		return Settings{}, err
	}
	return s, nil
}

// UpdateConfig sets the configuration settings in the config TOML file.
func UpdateConfig(s Settings) error {
	f, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	return toml.NewEncoder(w).Encode(s)
}
