package lex

type log struct {
	path    string
	source  string
	message string
}

type logs map[string][]log

func (l *logs) Add(path string, source string, message string) {
	logsMap := *l
	newLog := log{path: path, source: source, message: message}

	_, exists := logsMap[path]

	if !exists {
		logsMap[path] = []log{ newLog }
		return
	}

	logsMap[path] = append(logsMap[path], log{path: path, source: source, message: message})
}

func newLogs() *logs {
	return &logs{}
}
