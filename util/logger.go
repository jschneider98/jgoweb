package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var Logger *log.Logger

// factory method
func Log() *log.Logger {
	if Logger != nil {
		return Logger
	}

	Logger := log.New()

	Logger.SetFormatter(&log.JSONFormatter{})

	// @TEMP
	Logger.SetLevel(log.DebugLevel)

	return Logger
}

// defer Duration(Track("foo"))
func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

// debug duration to logger
func Duration(msg string, start time.Time) {
	Log().WithFields(
		log.Fields{
			"duration": fmt.Sprintf("%v", time.Since(start)),
		},
	).Debug(msg)
}
