package saslog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type F map[string]string

const (
	DefaultLevel = "INFO"
)

type Logger struct {
	l             *log.Logger
	level         string
	Data          F
	IncludePrefix bool
	prefix        string
	tmpData       F
}

func New(prefix string, defaultLevel string, data F) *Logger {
	l := new(Logger)
	l.l = log.New(l, "", 0)
	l.level = defaultLevel
	l.Data = data
	l.IncludePrefix = true
	l.prefix = prefix
	l.tmpData = F{}

	return l
}

func (logger *Logger) Info(msg string, data F) {
	logger.level = "INFO"
	logger.tmpData = data
	logger.l.Print(msg)
}
func (logger *Logger) Debug(msg string, data F) {
	logger.level = "DEBUG"
	logger.tmpData = data
	logger.l.Print(msg)
}

func (logger *Logger) Write(bytes []byte) (int, error) {
	output := time.Now().UTC().Format("2006-01-02 15:04:05.000")

	output += " " + logger.level

	if logger.IncludePrefix {
		output += " " + logger.prefix
	}

	// remove newline(s) added by log.Print
	output += " \"" + strings.TrimSpace(string(bytes)) + "\""

	for key, value := range logger.Data {
		output += " " + key + "=" + value
	}
	for key, value := range logger.tmpData {
		output += " " + key + "=" + value
	}

	// reset tmpData
	logger.tmpData = F{}

	w, e := fmt.Fprintln(os.Stderr, output)
	// reset level
	logger.level = DefaultLevel
	return w, e
}
