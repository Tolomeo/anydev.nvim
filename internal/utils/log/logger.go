package log

type logLevel string

const (
	Info  logLevel = "info"
	Warn  logLevel = "warn"
	Error logLevel = "error"
)

type log struct {
	Level   logLevel `json:"level" yaml:"level"`
	Message string   `json:"message" yaml:"message"`
}

type Logs map[string][]log

type Logger struct {
	key  string
	logs *Logs
}

func (l *Logger) Key() string {
	return l.key
}

func (l *Logger) SetKey(key string) {
	l.key = key
}

func (l *Logger) Log(newLog log) {
	logsMap := (*l.logs)

	_, exists := logsMap[l.key]

	if !exists {
		logsMap[l.key] = []log{newLog}
		return
	}

	logsMap[l.key] = append(logsMap[l.key], newLog)
}

func (l *Logger) Info(message string) {
	l.Log(log{Level: Info, Message: message})
}

func (l *Logger) Warn(message string) {
	l.Log(log{Level: Warn, Message: message})
}

func (l *Logger) Error(message string) {
	l.Log(log{Level: Error, Message: message})
}

func (l *Logger) Logs() *Logs {
	return l.logs
}

func NewLog(level logLevel, message string) log {
	return log{
		Level:   level,
		Message: message,
	}
}

func NewLogger(key string) *Logger {
	logs := Logs{}

	return &Logger{
		key:  key,
		logs: &logs,
	}
}
