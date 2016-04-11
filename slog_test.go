package slog

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

var config Config = Config{Prefix: "SAS:", Name: "sas"}

// Test that the standard logger has our prefix and a default level of 'INFO'
func TestLoggerStdLogOutput(t *testing.T) {
	msg := "test"
	buf := new(bytes.Buffer)
	config.Writer = buf

	l, err := New(config)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	log.SetFlags(0)
	log.SetOutput(l)

	// Use the standard logger
	log.Print(msg)

	s := strings.TrimSpace(buf.String())
	e := fmt.Sprintf("INFO %s \"%s\"", config.Name, msg)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}

// Test that the loggers output is formated properly and includes the default values
func TestLoggerOutput(t *testing.T) {
	// defaults
	key := "request_id"
	value := "1234"
	msg := "test"

	c := config

	buf := new(bytes.Buffer)
	c.Writer = buf

	c.SystemData = F{
		key: value,
	}
	c.AppData = F{
		key: value,
	}

	l, err := New(c)
	if err != nil {
		t.Error("Failed to create a Logger")
	}
	log.SetFlags(0)
	log.SetOutput(l)

	l.Info(msg, nil)

	s := strings.TrimSpace(buf.String())
	e := fmt.Sprintf("INFO %s \"%s\" %s=\"%s\"", config.Name, msg, key, value)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}

// Test that a new Logger has the correct defaults
func TestNewLoggerDefaults(t *testing.T) {

	l, err := New(config)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	// check Logger defaults match

	if l.prefix != config.Prefix {
		t.Error("prefix should be '", config.Prefix, "' but is '", l.prefix, "' instead")
	}

	if l.name != config.Name {
		t.Error("name should be '", config.Name, "' but is '", l.name, "' instead")
	}
}

// Test that subsequent logger calls use the correct level and doesn't use the previous value
func TestLoggerLevelReset(t *testing.T) {

	// defaults
	info_msg := "test info"
	debug_msg := "test debug"

	c := config
	buf := new(bytes.Buffer)
	c.Writer = buf

	l, err := New(c)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	log.SetFlags(0)
	log.SetOutput(l)

	// Send an INFO log entry
	l.Info(info_msg, nil)
	s := strings.TrimSpace(buf.String())
	buf.Reset()

	e := fmt.Sprintf("INFO %s \"%s\"", config.Name, info_msg)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}

	// Send a DEBUG log entry
	l.Debug(debug_msg, nil)
	s = strings.TrimSpace(buf.String())
	buf.Reset()
	e = fmt.Sprintf("%s %s \"%s\"", "DEBUG", config.Name, debug_msg)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}

	// Send a second INFO log entry
	l.Info(info_msg, nil)
	s = strings.TrimSpace(buf.String())
	buf.Reset()
	e = fmt.Sprintf("INFO %s \"%s\"", config.Name, info_msg)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}

// Test deriving a Logger from an existing Logger
func TestLoggerFromLogger(t *testing.T) {

	buf := new(bytes.Buffer)
	c := config
	c.Writer = buf
	l, err := New(c)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	c = Config{SystemData: F{"extra": "extra"}}

	nl := l.New(c)
	nl.Info("testing", nil)

	// Fields that should be inherited
	s, e := config.Prefix, nl.prefix
	if s != e {
		t.Errorf("'%s' != '%s'", s, e)
	}

	s, e = config.Name, nl.name
	if s != e {
		t.Errorf("'%s' != '%s'", s, e)
	}

	s = strings.TrimSpace(buf.String())
	e = fmt.Sprintf("INFO %s \"testing\" extra=\"extra\"", config.Name)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}

// Test that missing fields in config creates a default.
func TestLoggerConfigMissingFields(t *testing.T) {
	c := Config{}

	_, err := New(c)
	if err != nil {
		t.Error("New() shouldn't require any fields!")
	}

}

// Test that AppData and per-call data is prefixed with the Logger's prefix.
func TestLoggerAppDataPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	c := config
	c.Writer = buf
	c.SystemData = F{"system": "system"}
	c.AppData = F{"app": "app"}
	l, err := New(c)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	l.Info("testing", F{"per-call": "per-call"})

	s := strings.TrimSpace(buf.String())
	e := fmt.Sprintf("INFO %s \"testing\" system=\"system\" %s.app=\"app\" %s.per-call=\"per-call\"",
		config.Name, config.Prefix, config.Prefix)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}

// Test deriving a Logger from an existing Logger overwrites
// passed in config data.
func TestLoggerFromLoggerNewData(t *testing.T) {
	l, err := New(config)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	if l.prefix != config.Prefix {
		t.Errorf("%s != %s", l.prefix, config.Prefix)
	}

	if len(l.systemData) != len(config.SystemData) {
		t.Errorf("systemData has %d items, expected %d", len(l.systemData),
			len(config.SystemData))
	}

	if len(l.appData) != len(config.AppData) {
		t.Errorf("appData has %d items, expected %d", len(l.appData), len(config.AppData))
	}

	new_prefix := "new_prefix"
	new_name := "new_name"
	nl := l.New(Config{
		Prefix:     new_prefix,
		Name:       new_name,
		SystemData: F{"key": "value"},
		AppData:    F{"app_key": "app_value"},
	})

	// prefix shouldn't be changed
	if nl.prefix != l.prefix {
		t.Errorf("%s != %s", nl.prefix, l.prefix)
	}

	if nl.name != new_name {
		t.Errorf("%s != %s", nl.name, new_name)
	}

	if len(nl.systemData) != 1 {
		t.Errorf("systemData has %d items, expected %d", len(nl.systemData), 1)
	}

	if len(nl.appData) != 1 {
		t.Errorf("appData has %d items, expected %d", len(nl.appData), 1)
	}
}