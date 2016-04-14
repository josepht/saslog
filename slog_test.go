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

	l := New(config)

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

	c.SystemTags = T{
		key: value,
	}
	c.AppTags = T{
		key: value,
	}

	l := New(c)
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

	l := New(config)

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
	error_msg := "test error"

	c := config
	buf := new(bytes.Buffer)
	c.Writer = buf

	l := New(c)

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

	// Send an ERROR log entry
	l.Error(error_msg, nil)
	s = strings.TrimSpace(buf.String())
	buf.Reset()
	e = fmt.Sprintf("%s %s \"%s\"", "ERROR", config.Name, error_msg)
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
	l := New(c)

	c = Config{SystemTags: T{"extra": "extra"}}

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

// Test that AppTags and per-call tags are prefixed with the Logger's prefix.
func TestLoggerAppTagsPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	c := config
	c.Writer = buf
	c.SystemTags = T{"system": "system"}
	c.AppTags = T{"app": "app"}
	l := New(c)

	l.Info("testing", T{"per-call": "per-call"})

	s := strings.TrimSpace(buf.String())
	e := fmt.Sprintf("INFO %s \"testing\" system=\"system\" %s.app=\"app\" %s.per-call=\"per-call\"",
		config.Name, config.Prefix, config.Prefix)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}

// Test deriving a Logger from an existing Logger overwrites
// passed in config data except tagss.
func TestLoggerFromLoggerNewTags(t *testing.T) {
	c := config
	c.SystemTags = T{
		"orig_key": "orig_value",
	}
	c.AppTags = T{
		"app_orig_key": "app_orig_value",
	}
	l := New(c)

	if l.prefix != c.Prefix {
		t.Errorf("%s != %s", l.prefix, c.Prefix)
	}

	if len(l.systemTags) != len(c.SystemTags) {
		t.Errorf("systemTags has %d items, expected %d", len(l.systemTags),
			len(c.SystemTags))
	}

	if len(l.appTags) != len(c.AppTags) {
		t.Errorf("appTags has %d items, expected %d", len(l.appTags), len(c.AppTags))
	}

	new_prefix := "new_prefix"
	new_name := "new_name"
	nl := l.New(Config{
		Prefix:     new_prefix,
		Name:       new_name,
		SystemTags: T{"orig_key": "new_value"},
		AppTags:    T{"app_orig_key": "app_new_value"},
	})

	// prefix shouldn't be changed
	if nl.prefix != l.prefix {
		t.Errorf("%s != %s", nl.prefix, l.prefix)
	}

	if nl.name != new_name {
		t.Errorf("%s != %s", nl.name, new_name)
	}

	// Test that the new logger hasn't modified the original tag values
	if len(nl.systemTags) != 1 {
		t.Errorf("systemTags has %d items, expected %d", len(nl.systemTags), 1)
	}

	if len(nl.appTags) != 1 {
		t.Errorf("appTags has %d items, expected %d", len(nl.appTags), 1)
	}

	if nl.systemTags["orig_key"] != "orig_value" {
		t.Errorf("systemTags has %s for 'orig_key', expected %s", nl.systemTags["orig_key"],
			"orig_value")
	}

	if nl.appTags["app_orig_key"] != "app_orig_value" {
		t.Errorf("appTags has %s for 'app_orig_key', expected %s", nl.appTags["app_orig_key"],
			"app_orig_value")
	}

	// Test that the original logger wasn't modified
	if len(l.systemTags) != 1 {
		t.Errorf("systemTags has %d items, expected %d", len(l.systemTags), 1)
	}

	if len(l.appTags) != 1 {
		t.Errorf("appTags has %d items, expected %d", len(l.appTags), 1)
	}

	if l.systemTags["orig_key"] != "orig_value" {
		t.Errorf("systemTags has %s for 'orig_key', expected %s", l.systemTags["orig_key"],
			"orig_value")
	}

	if l.appTags["app_orig_key"] != "app_orig_value" {
		t.Errorf("appTags has %s for 'app_orig_key', expected %s", l.appTags["app_orig_key"],
			"app_orig_value")
	}
}

// Test that a nil Logger doesn't blow up when its methods are called.
func TestNilLogger(t *testing.T) {
	var l *Logger

	l.Info("test", nil)
	l.Debug("test", nil)
	l.Error("test", nil)

	l.Write([]byte("test"))

	nl := l.New(Config{})

	if nl != nil {
		t.Error("Creating a Logger from a nil Logger should yeild nil")
	}
}

// Test that an empty prefix doesn't include '.' in the key/value pairs
func TestEmptyPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	c := config
	c.Writer = buf
	c.Prefix = ""
	l := New(c)

	l.Info("testing", T{"per-call": "per-call"})

	s := strings.TrimSpace(buf.String())
	e := fmt.Sprintf("INFO %s \"testing\" per-call=\"per-call\"", config.Name)
	if !strings.Contains(s, e) {
		t.Errorf("'%s' doesn't contain '%s'", s, e)
	}
}
