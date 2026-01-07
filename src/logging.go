package main

import (
	"log"
	"os"
	"strings"
)

type logLevel int

const (
	logLevelInfo logLevel = iota
	logLevelDebug
)

var currentLogLevel logLevel

func initLogLevel() {
	level := strings.ToUpper(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	switch level {
	case "DEBUG":
		currentLogLevel = logLevelDebug
	default:
		currentLogLevel = logLevelInfo
	}
}

func logInfo(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func logDebug(format string, v ...interface{}) {
	if currentLogLevel >= logLevelDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func logError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
