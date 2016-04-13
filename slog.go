package slog

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
}

type Config struct {
	Writer     io.Writer // optional
	Name       string    // optional
	Prefix     string    // optional
	SystemData F
	AppData    F
}

// Create a new logger based on the passed in config.
func New(c Config) *Logger {
	l := new(Logger)

	if c.Writer == nil {
		c.Writer = os.Stderr
	}
	l.l = log.New(c.Writer, "", 0)

	if c.Name == "" {
		c.Name = "-"
	}

	l.name = c.Name
	l.prefix = c.Prefix
	l.systemData = c.SystemData
	l.appData = c.AppData

	return l
}

// Create a new logger based on the current logger.  Any
// config values will overwrite the values for the current
// logger.
func (l *Logger) New(c Config) *Logger {
	if l == nil || l.l == nil {
		return nil
	}
	nl := new(Logger)
	nl.l = l.l

	// copy original Logger fields
	nl.prefix = l.prefix
	nl.name = l.name
	nl.systemData = F{}
	nl.appData = F{}

	if c.Name != "" {
		nl.name = c.Name
	}

	for k, v := range l.systemData {
		nl.systemData[k] = v
	}

	for k, v := range l.appData {
		nl.appData[k] = v
	}

	if c.SystemData != nil {
		// Create the system data map if it doesn't exist.
		if nl.systemData == nil {
			nl.systemData = F{}
		}
		for k, v := range c.SystemData {
			// Don't update existing fields
			if _, ok := nl.systemData[k]; !ok {
				nl.systemData[k] = v
			}
		}
	}

	if c.AppData != nil {
		// Create the app data map if it doesn't exist.
		if nl.appData == nil {
			nl.appData = F{}
		}
		for k, v := range c.AppData {
			// Don't update existing fields
			if _, ok := nl.appData[k]; !ok {
				nl.appData[k] = v
			}
		}
	}

	return nl
}

func (l *Logger) log(msg string, level string, data F) {
	if l == nil || l.l == nil {
		return
	}
	ts := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	d := ""

	// Add the system data
	for key, value := range l.systemData {
		d += fmt.Sprintf(" %s=%s", key, strconv.Quote(value))
	}

	prefix := ""
	if l.prefix == "" {
		prefix = l.prefix
	} else {
		prefix = fmt.Sprintf("%s.", l.prefix)
	}

	// Add the application data
	for key, value := range l.appData {
		d += fmt.Sprintf(" %s%s=%s", prefix, key, strconv.Quote(value))
	}

	// Add the per-call data
	for key, value := range data {
		d += fmt.Sprintf(" %s%s=%s", prefix, key, strconv.Quote(value))
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

func (l *Logger) Error(msg string, data F) {
	l.log(msg, "ERROR", data)
}

func (l *Logger) Write(bytes []byte) (int, error) {
	// remove trailing whitespace
	l.Info(strings.TrimSpace(string(bytes)), nil)
	return len(bytes), nil
}
