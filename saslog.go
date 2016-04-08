package saslog

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type F map[string]string

const (
	DefaultLevel = "INFO"
)

type Logger struct {
	l          *log.Logger
	systemData F
	appData    F
	name       string
}

func New(w io.Writer, name string, sd F, ad F) *Logger {
	l := new(Logger)
	l.l = log.New(l, "", 0)
	l.l.SetFlags(0)
	l.l.SetOutput(w)
	l.name = name
	l.systemData = sd
	l.appData = ad

	return l
}

func (l *Logger) log(msg string, level string, data F) {
	ts := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	d := ""

	// Add the system data
	for key, value := range l.systemData {
		d += fmt.Sprintf(" %s=%s", key, value)
	}

	var s string
	if l.systemData["service"] == "" {
		s = l.name
	} else {
		s = l.systemData["service"]
	}

	// Add the application data
	for key, value := range l.appData {
		d += fmt.Sprintf(" %s.%s=%s", s, key, value)
	}

	// Add the per-call data
	for key, value := range data {
		d += fmt.Sprintf(" %s.%s=%s", s, key, value)
	}

	output := fmt.Sprintf("%s %s %s \"%s\"%s", ts, level, l.name, msg, d)

	l.l.Print(output)

}

func (l *Logger) Info(msg string, data F) {
	l.log(msg, "INFO", data)
}
func (l *Logger) Debug(msg string, data F) {
	l.log(msg, "DEBUG", data)
}

func (l *Logger) Write(bytes []byte) (int, error) {
	// remove trailing whitespace
	l.Info(strings.TrimSpace(string(bytes)), nil)
	return len(bytes), nil
}
