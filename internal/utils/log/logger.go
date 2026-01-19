package log

import (
	"fmt"
	stdLog "log"
	"os"
)

type logLevel string

const (
	Info  logLevel = "INFO"
	Warn  logLevel = "WARN"
	Error logLevel = "ERROR"
)

type log struct {
	Level   logLevel `json:"level" yaml:"level"`
	Message string   `json:"message" yaml:"message"`
}

type Logger struct {
	stdOut *stdLog.Logger
}

func (l *Logger) Log(newLog log) {
	l.stdOut.Printf("[%s] %s\n", newLog.Level, newLog.Message)
}

func (l *Logger) Info(message string) {
	l.Log(log{Level: Info, Message: message})
}

func (l *Logger) Infof(message string, args ...any) {
	l.Log(log{Level: Info, Message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Warn(message string) {
	l.Log(log{Level: Warn, Message: message})
}

func (l *Logger) Warnf(message string, args ...any) {
	l.Log(log{Level: Warn, Message: fmt.Sprintf(message, args...)})
}

func (l *Logger) Error(message string) {
	l.Log(log{Level: Error, Message: message})
}

func (l *Logger) Errorf(message string, args ...any) {
	l.Log(log{Level: Error, Message: fmt.Sprintf(message, args...)})
}

func NewLogger() *Logger {
	return &Logger{
		stdOut: stdLog.New(os.Stdout, "", stdLog.Ldate|stdLog.Ltime),
	}
}
