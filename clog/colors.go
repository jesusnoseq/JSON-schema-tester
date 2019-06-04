package clog

import (
	"github.com/fatih/color"
)

const sufix = "\n"

var error = color.New(color.FgRed).PrintfFunc()
var success = color.New(color.FgGreen).PrintfFunc()
var info = color.New(color.FgCyan).PrintfFunc()
var debug = color.New().PrintfFunc()

// Success print success msg
func Success(format string, a ...interface{}) {
	success(format+sufix, a...)
}

// Error print error msg
func Error(format string, a ...interface{}) {
	error(format+sufix, a...)
}

// Info print info msg
func Info(format string, a ...interface{}) {
	info(format+sufix, a...)
}

// Debug print debug msg
func Debug(format string, a ...interface{}) {
	debug(format+sufix, a...)
}
