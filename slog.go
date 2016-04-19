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

type T map[string]string

type Logger struct {
	l          *log.Logger
	systemTags T
	appTags    T
	name       string
	prefix     string
}

type RootConfig struct {
	Writer io.Writer // optional
	Config           // optional
}

type Config struct {
	Name       string // optional
	Prefix     string // optional
	SystemTags T
	AppTags    T
}

// Create a new logger based on the passed in config.
func New(c RootConfig) *Logger {
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
	l.systemTags = c.SystemTags
	l.appTags = c.AppTags

	return l
}

// Create a new logger based on the current logger.  Any
// config values will overwrite the values for the current
// logger.
func (l *Logger) New(c Config) *Logger {
	if l == nil {
		return nil
	}
	nl := new(Logger)
	nl.l = l.l

	// copy original Logger fields
	nl.prefix = l.prefix
	nl.name = l.name
	nl.systemTags = T{}
	nl.appTags = T{}

	if c.Name != "" {
		nl.name = c.Name
	}

	for k, v := range l.systemTags {
		nl.systemTags[k] = v
	}

	for k, v := range l.appTags {
		nl.appTags[k] = v
	}

	if c.SystemTags != nil {
		for k, v := range c.SystemTags {
			// Don't update existing fields
			if _, ok := nl.systemTags[k]; !ok {
				nl.systemTags[k] = v
			}
		}
	}

	if c.AppTags != nil {
		for k, v := range c.AppTags {
			// Don't update existing fields
			if _, ok := nl.appTags[k]; !ok {
				nl.appTags[k] = v
			}
		}
	}

	return nl
}

func (l *Logger) log(msg string, level string, tags T) {
	if l == nil {
		return
	}
	ts := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	d := ""

	// Add the system tags
	for key, value := range l.systemTags {
		d += fmt.Sprintf(" %s=%s", key, strconv.Quote(value))
	}

	prefix := ""
	if l.prefix == "" {
		prefix = l.prefix
	} else {
		prefix = fmt.Sprintf("%s.", l.prefix)
	}

	// Add the application tags
	for key, value := range l.appTags {
		d += fmt.Sprintf(" %s%s=%s", prefix, key, strconv.Quote(value))
	}

	// Add the per-call tags
	for key, value := range tags {
		d += fmt.Sprintf(" %s%s=%s", prefix, key, strconv.Quote(value))
	}

	output := fmt.Sprintf("%s %s %s %s%s", ts, level, l.name,
		strconv.Quote(msg), d)

	l.l.Print(output)

}

func (l *Logger) Info(msg string, tags T) {
	l.log(msg, "INFO", tags)
}

func (l *Logger) Debug(msg string, tags T) {
	l.log(msg, "DEBUG", tags)
}

func (l *Logger) Error(msg string, tags T) {
	l.log(msg, "ERROR", tags)
}

func (l *Logger) Write(bytes []byte) (int, error) {
	// remove trailing whitespace
	l.Info(strings.TrimSpace(string(bytes)), nil)
	return len(bytes), nil
}
