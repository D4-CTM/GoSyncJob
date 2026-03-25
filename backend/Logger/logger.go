package logger

import (
	"fmt"
	"log"
)

type Level string

const (
	DEBUG   Level = "[DEBUG]"
	INFO    Level = "[INFO]"
	WARNING Level = "[WARN]"
	ERROR   Level = "[ERROR]"
	FATAL   Level = "[FATAL]"
)

func LogMessage(level Level, msg string, args ...any) {
	log.Printf(string(level) + " " + msg, args...)
}

func LogDebug(msg string, args ...any) {
	LogMessage(DEBUG, msg, args...)
}

func LogInfo(msg string, args ...any) {
	LogMessage(INFO, msg, args...)
}

func LogWarn(msg string, args ...any) {
	LogMessage(WARNING, msg, args...)
}

func LogErr(msg string, args ...any) error {
	out := fmt.Sprintf(msg, args...);
	LogMessage(ERROR, out)
	return fmt.Errorf("%s %s", ERROR, out)
}

func LogFatal(msg string, args ...any) {
	log.Fatalf(string(FATAL) + " " + msg, args...)
}
