package clog

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var red = color.New(color.FgRed).SprintfFunc()
var green = color.New(color.FgGreen).SprintfFunc()
var cyan = color.New(color.FgCyan).SprintfFunc()
var white = color.New(color.FgHiWhite).SprintfFunc()

var logger = initLogger()

// ErrorCounter counter for errors
var errorCounter = 0

func initLogger() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})
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
