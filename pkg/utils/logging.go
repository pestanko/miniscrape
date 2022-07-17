package utils

import (
	"io"
	"os"
	"path"

	"github.com/pestanko/miniscrape/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// MakeAccessLog creates an instance of the access logger
func MakeAccessLog(logCfg *config.LogConfig) zerolog.Logger {
	return makeLogger(logCfg, "access.log")
}

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(logCfg *config.LogConfig) {
	mw := makeLogWritters(logCfg, "app.log")
	log.Logger = log.Output(mw).
		With().
		Str("compnent", "app").
		Logger()
}

func makeLogger(logCfg *config.LogConfig, file string) zerolog.Logger {

	mw := makeLogWritters(logCfg, file)

	logger := zerolog.
		New(mw).
		With().
		Timestamp().
		Logger()

	return logger
}

func makeLogWritters(logCfg *config.LogConfig, file string) io.Writer {
	var writers []io.Writer

	if logCfg.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		writers = append(writers, os.Stderr)
	}

	if logCfg.Dir != "" {
		writers = append(writers, newRollingFile(logCfg, file))
	}

	return io.MultiWriter(writers...)
}

func newRollingFile(logCfg *config.LogConfig, file string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(logCfg.Dir, file),
		MaxBackups: 3,  // files
		MaxSize:    10, // megabytes
		MaxAge:     30, // days
	}
}
