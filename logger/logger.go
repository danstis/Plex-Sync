package logger

import (
	"io"
	"os"

	"github.com/danstis/Plex-Sync/config"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// CreateLogger builds a file logger using Lumberjack
func CreateLogger(filename string) io.Writer {
	return io.MultiWriter(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.Settings.General.MaxLogSize, // megabytes
		MaxBackups: config.Settings.General.MaxLogCount,
		MaxAge:     config.Settings.General.MaxLogAge, //days
	}, os.Stdout)
}
