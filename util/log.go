package util

import (
	log "github.com/sirupsen/logrus"
)

// SetupLog sets the log level and log format
func SetupLog(logLevel string) {
	setLogLevel(logLevel)

	// Set log format
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{log.FieldKeyMsg: "message"},
	})
}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
