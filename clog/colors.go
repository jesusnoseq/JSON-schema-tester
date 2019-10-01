package clog

import (
	"github.com/fatih/color"
	"github.com/jesusnoseq/JSON-schema-tester/config"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var red = color.New(color.FgRed).SprintfFunc()
var green = color.New(color.FgGreen).SprintfFunc()
var cyan = color.New(color.FgCyan).SprintfFunc()
var white = color.New(color.FgHiWhite).SprintfFunc()

var logger *log.Logger

// ErrorCounter counter for errors
var errorCounter = 0

// Warn counter for errors
var warnCounter = 0

var warnsAllowed = 0

func InitLogger(conf config.PathConfig) *log.Logger {
	warnsAllowed = conf.WarnsAllowed
	logLevel, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatal("Log level " + conf.LogLevel + " is not valid")
	}

	l := log.New()
	l.SetLevel(logLevel)
	l.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
	logger = l
	return l
}

// Success print success msg
func Success(format string, a ...interface{}) {
	logger.Debug(green(format, a...))
}

// Error print error msg
func Error(format string, args ...interface{}) {
	errorCounter++
	logger.Error(red(strconv.Itoa(errorCounter)+" "+format, args...))
}

// Error print error msg
func Warn(format string, args ...interface{}) {
	warnCounter++
	logger.Warn(red(strconv.Itoa(warnCounter)+" "+format, args...))
	if warnCounter > warnsAllowed {
		log.Errorf("Wow! more than %d warnings detected, please check latest added schemas", warnsAllowed)
	}
}

// Info print info msg
func Info(format string, a ...interface{}) {
	logger.Info(cyan(format, a...))
}

// Debug print debug msg
func Debug(format string, a ...interface{}) {
	logger.Debug(white(format, a...))
}

func GetErrorsPrinted() int {
	return errorCounter
}
