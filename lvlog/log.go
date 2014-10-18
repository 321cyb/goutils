// Package log implements a wrapper around the Go standard library's
// logging package. Clients should set the current log level; only
// messages below that level will actually be logged. For example, if
// Level is set to LevelWarning, only log messages at the Warning,
// Error, and Critical levels will be logged.

// TODO: What if No Init function is called?

package lvlog

import (
	"errors"
	"flag"
	"fmt"
	golog "log"
	"os"
)

// The following constants represent logging levels in increasing levels of seriousness.
const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
)

var levelPrefix = [...]string{
	LevelDebug:   "[DEBUG] ",
	LevelInfo:    "[INFO] ",
	LevelWarning: "[WARNING] ",
	LevelError:   "[ERROR] ",
}

// Level stores the current logging level.
var level = LevelInfo
var logFile = os.Stderr

var levelString = "INFO"
var filename = ""
var logger golog.Logger

func init() {
	flag.StringVar(&levelString, "log-level", "INFO", "log level, can only be DEBUG, INFO, WARNING, ERROR")
	flag.StringVar(&filename, "log-file", "", "log file name")
}

func setFileName(name string) error {
	var err error
	if len(name) > 0 {
		logFile, err = os.OpenFile(name, os.O_RDWR, 0666)
	}
	golog.SetOutput(logFile)
	golog.SetFlags(golog.LstdFlags | golog.Lshortfile)
	return err
}

// InitFromArgs is the preferred way to init this library.
func InitFromArgs() {
	switch levelString {
	case "DEBUG":
		level = LevelDebug
	case "INFO":
		level = LevelInfo
	case "WARNING":
		level = LevelWarning
	case "ERROR":
		level = LevelError
	default:
		fmt.Fprintln(os.Stderr, "log-level can only be DEBUG, INFO, WARNING, ERROR")
		os.Exit(1)
	}
	setFileName(filename)
}

// InitLevelAndFile is used to manually init this library
func InitLevelAndFile(l int, p string) error {
	if l <= LevelError && l >= LevelDebug {
		level = l
	} else {
		return errors.New("wrong level number")
	}
	return setFileName(p)
}

func outputf(l int, format string, v []interface{}) {
	if l >= level {
		golog.Printf(fmt.Sprint(levelPrefix[l], format), v...)
	}
}

func output(l int, v []interface{}) {
	if l >= level {
		golog.Print(levelPrefix[l], fmt.Sprint(v...))
	}
}

// Fatalf logs a formatted message at the "critical" level. The
// arguments are handled in the same manner as fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	outputf(LevelError, format, v)
	os.Exit(1)
}

// Fatal logs its arguments at the "critical" level.
func Fatal(v ...interface{}) {
	output(LevelError, v)
	os.Exit(1)
}

// Errorf logs a formatted message at the "error" level. The arguments
// are handled in the same manner as fmt.Printf.
func Errorf(format string, v ...interface{}) {
	outputf(LevelError, format, v)
}

// Error logs its arguments at the "error" level.
func Error(v ...interface{}) {
	output(LevelError, v)
}

// Warningf logs a formatted message at the "warning" level. The
// arguments are handled in the same manner as fmt.Printf.
func Warningf(format string, v ...interface{}) {
	outputf(LevelWarning, format, v)
}

// Warning logs its arguments at the "warning" level.
func Warning(v ...interface{}) {
	output(LevelWarning, v)
}

// Infof logs a formatted message at the "info" level. The arguments
// are handled in the same manner as fmt.Printf.
func Infof(format string, v ...interface{}) {
	outputf(LevelInfo, format, v)
}

// Info logs its arguments at the "info" level.
func Info(v ...interface{}) {
	output(LevelInfo, v)
}

// Debugf logs a formatted message at the "debug" level. The arguments
// are handled in the same manner as fmt.Printf.
func Debugf(format string, v ...interface{}) {
	outputf(LevelDebug, format, v)
}

// Debug logs its arguments at the "debug" level.
func Debug(v ...interface{}) {
	output(LevelDebug, v)
}
