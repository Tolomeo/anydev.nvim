package log

import (
	"fmt"
	stdLog "log"
	"os"
	"slices"
)

type level string

const (
	error   level = "ERROR"
	warn    level = "WARN"
	info    level = "INFO"
	verbose level = "VERBOSE"
	silly   level = "SILLY"
	debug   level = "DEBUG"
)

type color string

const (
	reset   color = "\033[0m"
	red     color = "\033[31m"
	green   color = "\033[32m"
	yellow  color = "\033[33m"
	blue    color = "\033[34m"
	magenta color = "\033[35m"
	cyan    color = "\033[36m"
	gray    color = "\033[37m"
	// white   color = "\033[97m"
)

var levels = [5]level{error, warn, info, verbose, silly}
var colors = map[level]color{
	error:   red,
	warn:    yellow,
	info:    green,
	verbose: blue,
	silly:   gray,
	debug:   cyan,
}

type log struct {
	level   level
	message string
}

type Logger struct {
	level  uint
	prefix string
	stdOut *stdLog.Logger
}

func (l *Logger) log(newLog log) {
	newLogLevel := slices.Index(levels[:], newLog.level)

	if newLogLevel > int(l.level) {
		return
	}

	level := fmt.Sprintf("%s%s%s", colors[newLog.level], newLog.level, reset)
	prefix := fmt.Sprintf("[%s%s%s]", magenta, l.prefix, reset)
	message := newLog.message

	l.stdOut.Printf("%s%s: %s\n", level, prefix, message)
}

func (l *Logger) Info(message string) {
	l.log(log{level: info, message: message})
}

func (l *Logger) Infof(message string, args ...any) {
	l.log(log{level: info, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Warn(message string) {
	l.log(log{level: warn, message: message})
}

func (l *Logger) Warnf(message string, args ...any) {
	l.log(log{level: warn, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Error(message string) {
	l.log(log{level: error, message: message})
}

func (l *Logger) Errorf(message string, args ...any) {
	l.log(log{level: error, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Verbose(message string) {
	l.log(log{level: verbose, message: message})
}

func (l *Logger) Verbosef(message string, args ...any) {
	l.log(log{level: verbose, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Silly(message string) {
	l.log(log{level: silly, message: message})
}

func (l *Logger) Sillyf(message string, args ...any) {
	l.log(log{level: silly, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Debug(message string) {
	message = fmt.Sprintf("\n****************************\n%s\n****************************\n", message)

	l.log(log{level: debug, message: message})
}

func (l *Logger) Debugf(message string, args ...any) {
	message = fmt.Sprintf("\n****************************\n%s\n****************************\n", fmt.Sprintf(message, args...))

	l.log(log{level: debug, message: message})
}

func NewLogger(prefix string, level uint) *Logger {
	return &Logger{
		level:  level,
		prefix: prefix,
		stdOut: stdLog.New(os.Stdout, "", stdLog.Ldate|stdLog.Ltime),
	}
}
