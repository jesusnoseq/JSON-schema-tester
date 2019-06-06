package clog

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var error = color.New(color.FgRed).PrintfFunc()
var success = color.New(color.FgGreen).PrintfFunc()
var info = color.New(color.FgCyan).PrintfFunc()
var debug = color.New().PrintfFunc()

var logger = initLogger()
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
	//logrus.Debug()
	//success(format+"\n", a...)
	f := color.New(color.FgGreen).SprintfFunc()
	logger.Debug(f(format, a...))
}

// Error print error msg
func Error(format string, args ...interface{}) {
	errorCounter++
	error(strconv.Itoa(errorCounter)+" "+format+"\n", args...)
	//logger.Error(args...)
}

// Info print info msg
func Info(format string, a ...interface{}) {
	info(format, a...)
}

// Debug print debug msg
func Debug(format string, a ...interface{}) {
	debug(format, a...)
}
