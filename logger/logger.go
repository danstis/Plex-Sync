package logger

import (
	"io"
	"os"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// CreateLogger builds a file logger using Lumberjack
func CreateLogger(filename string, size, count, age int) io.Writer {
	return io.MultiWriter(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    size, // megabytes
		MaxBackups: count,
		MaxAge:     age, //days
	}, os.Stdout)
}
