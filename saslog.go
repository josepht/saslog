package saslog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type F map[string]string

type Logger struct {
	l          *log.Logger
	systemData F
	appData    F
	name       string
	prefix     string
	writer     io.Writer
}

type Config struct {
	Writer     io.Writer // optional
	Name       string    // optional
	Prefix     string    // optional
	SystemData F
	AppData    F
}

// Create a new logger based on the passed in config.
func New(c Config) (*Logger, error) {
	l := new(Logger)

	if c.Writer == nil {
		c.Writer = os.Stderr
	}
	l.l = log.New(c.Writer, "", 0)
	l.writer = c.Writer

	if c.Name == "" {
		c.Name = "-"
	}

	l.name = c.Name
	l.prefix = c.Prefix
	l.systemData = c.SystemData
	l.appData = c.AppData

	return l, nil
}

// Create a new logger based on the current logger.  Any
// config values will overwrite the values for the current
// logger.
func (l *Logger) New(c Config) *Logger {
	nl := new(Logger)
	nl.l = log.New(l, "", 0)

	// copy original Logger fields
	nl.prefix = l.prefix
	nl.name = l.name
	nl.systemData = l.systemData
	nl.appData = l.appData

	if c.Writer == nil {
		c.Writer = l.writer
	}
	nl.l.SetOutput(c.Writer)
	nl.writer = c.Writer

	if c.Name != "" {
		nl.name = c.Name
	}

	// TODO: combine old and new data
	if c.SystemData != nil {
		nl.systemData = c.SystemData
	}

	// TODO: combine old and new data
	if c.AppData != nil {
		nl.appData = c.AppData
	}

	return nl
}

func (l *Logger) log(msg string, level string, data F) {
	ts := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	d := ""

	// Add the system data
	for key, value := range l.systemData {
		d += fmt.Sprintf(" %s=%s", key, strconv.Quote(value))
	}

	// Add the application data
	for key, value := range l.appData {
		d += fmt.Sprintf(" %s.%s=%s", l.prefix, key, strconv.Quote(value))
	}

	// Add the per-call data
	for key, value := range data {
		d += fmt.Sprintf(" %s.%s=%s", l.prefix, key, strconv.Quote(value))
	}

	output := fmt.Sprintf("%s %s %s %s%s", ts, level, l.name,
		strconv.Quote(msg), d)

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
