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
	debug   level = "DEBUG"
)

var levels = [4]level{error, warn, info, verbose}

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

	l.stdOut.Printf("%s%s: %s\n", newLog.level, l.prefix, newLog.message)
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

func (l *Logger) Debug(message string) {
	message = fmt.Sprintf("\n****************************\n%s\n****************************\n", message)

	l.log(log{level: debug, message: message})
}

func (l *Logger) Debugf(message string, args ...any) {
	message = fmt.Sprintf("\n****************************\n%s\n****************************\n", fmt.Sprintf(message, args...))

	l.log(log{level: debug, message: message})
}

func NewLogger(prefix string) *Logger {
	if prefix != "" {
		prefix = fmt.Sprintf("[%s]", prefix)
	}

	return &Logger{
		level:  2,
		prefix: prefix,
		stdOut: stdLog.New(os.Stdout, "", stdLog.Ldate|stdLog.Ltime),
	}
}
