package applog

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig logger configuration
type LogConfig struct {
	// Dir where to store log files
	Dir string `json:"dir"`
	// ConsoleLoggingEnabled whether logger should use console logging
	ConsoleLoggingEnabled bool `json:"console_logging_enabled"`
}

// MakeAccessLog creates an instance of the access logger
func MakeAccessLog(logCfg *LogConfig) zerolog.Logger {
	return makeLogger(logCfg, "access.log").With().
		Str("logger_kind", "access").
		Logger()
}

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(logCfg *LogConfig) {
	mw := makeLogWriters(logCfg, "app.log")
	log.Logger = log.Output(mw).
		With().
		Timestamp().
		Str("compnent", "app").
		Logger()
	zerolog.DefaultContextLogger = &log.Logger
}

func makeLogger(logCfg *LogConfig, file string) zerolog.Logger {

	mw := makeLogWriters(logCfg, file)

	logger := zerolog.
		New(mw).
		With().
		Timestamp().
		Logger()

	return logger
}

func makeLogWriters(logCfg *LogConfig, file string) io.Writer {
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

func newRollingFile(logCfg *LogConfig, file string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(logCfg.Dir, file),
		MaxBackups: 3,  // files
		MaxSize:    10, // megabytes
		MaxAge:     30, // days
	}
}
