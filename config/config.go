package config

import (
	"time"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/spf13/viper"
)

// Settings exports the defined application settings
var Settings Environment

// Environment contains the environment configuration
type Environment struct {
	General      generalSettings
	LocalServer  plex.Host
	RemoteServer plex.Host
	Webui        webuiSettings
}

type generalSettings struct {
	SyncInterval     time.Duration
	WebserverPort    int
	Logfile          string
	WebserverLogfile string
	MaxLogSize       int
	MaxLogCount      int
	MaxLogAge        int
}

type webuiSettings struct {
	CacheLifetime int
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Settings.General.SyncInterval = viper.GetDuration("general.interval") * time.Second
	Settings.General.WebserverPort = viper.GetInt("general.webserverport")
	Settings.General.WebserverLogfile = viper.GetString("general.webserverlogfile")
	Settings.General.Logfile = viper.GetString("general.logfile")
	Settings.General.MaxLogSize = viper.GetInt("general.maxlogsize")
	Settings.General.MaxLogCount = viper.GetInt("general.maxlogcount")
	Settings.General.MaxLogAge = viper.GetInt("general.maxlogage")
	Settings.LocalServer = plex.Host{
		Name:     viper.GetString("localServer.name"),
		Hostname: viper.GetString("localServer.hostname"),
		Port:     viper.GetInt("localServer.port"),
		Ssl:      viper.GetBool("localServer.usessl"),
	}
	Settings.RemoteServer = plex.Host{
		Name:     viper.GetString("remoteServer.name"),
		Hostname: viper.GetString("remoteServer.hostname"),
		Port:     viper.GetInt("remoteServer.port"),
		Ssl:      viper.GetBool("remoteServer.usessl"),
	}
	Settings.Webui.CacheLifetime = viper.GetInt("webui.cacheLifetime") * -1
}
