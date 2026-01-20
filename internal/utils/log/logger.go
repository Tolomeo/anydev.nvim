package log

import (
	"fmt"
	stdLog "log"
	"os"
)

type level string

const (
	info  level = "INFO"
	warn  level = "WARN"
	error level = "ERROR"
)

type log struct {
	level   level
	message string
}

type Logger struct {
	stdOut *stdLog.Logger
}

func (l *Logger) Log(newLog log) {
	l.stdOut.Printf("[%s] %s\n", newLog.level, newLog.message)
}

func (l *Logger) Info(message string) {
	l.Log(log{level: info, message: message})
}

func (l *Logger) Infof(message string, args ...any) {
	l.Log(log{level: info, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Warn(message string) {
	l.Log(log{level: warn, message: message})
}

func (l *Logger) Warnf(message string, args ...any) {
	l.Log(log{level: warn, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Error(message string) {
	l.Log(log{level: error, message: message})
}

func (l *Logger) Errorf(message string, args ...any) {
	l.Log(log{level: error, message: fmt.Sprintf(message, args...)})
}

func (l *Logger) SetPrefix(prefix string) {
	l.stdOut.SetPrefix(prefix)
}

func NewLogger() *Logger {
	return &Logger{
		stdOut: stdLog.New(os.Stdout, "", stdLog.Ldate|stdLog.Ltime),
	}
}
