package util

import (
	"time"
	"fmt"
	"strconv"
	"os"
)

var Debug bool

//
func Debugln(a ...interface{}) (int, error) {
	
	if Debugging() {
		return fmt.Println(a...)
	}

	return 0, nil
}

//
func Debugf(format string, a ...interface{}) (int, error) {
	
	if Debugging() {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}

// Use defer to test method executuion time
func DebugTimeTrack(start time.Time, name string) {
	
	if !Debugging() {
		return
	}

	elapsed := time.Since(start)
	Debugf("%s took %s\n", name, elapsed)
}

//
func SetDebug(debug bool) {
	Debug = debug
}

//
func Debugging() bool {

	if Debug {
		return true
	}

	logLevelStr := os.Getenv("JGOLOGLEVEL")

	if logLevelStr == "" {
		return false
	}

	 logLevel, err := strconv.Atoi(logLevelStr)

	 if err != nil {
	 	return false
	 }

	 return logLevel <= 2
}
