package saslog

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

// Test that the standard logger has our prefix and a default level of 'INFO'
func TestLoggerStdLogOutput(t *testing.T) {
	prefix := "SAS:"
	data := F{}
	app_data := F{}
	msg := "test"

	buf := new(bytes.Buffer)

	l := New(buf, prefix, data, app_data)
	log.SetFlags(0)
	log.SetOutput(l)

	// Use the standard logger
	log.Print(msg)

	strout := strings.TrimSpace(buf.String())

	expected := fmt.Sprintf("INFO %s \"%s\"", prefix, msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

}

// Test that the loggers output is formated properly and includes the default values
func TestLoggerOutput(t *testing.T) {
	// defaults
	prefix := "SAS:"
	level := "INFO"
	key := "request_id"
	value := "1234"
	data := F{
		key: value,
	}
	app_data := F{
		key: value,
	}
	msg := "test"
	buf := new(bytes.Buffer)

	l := New(buf, prefix, data, app_data)
	log.SetFlags(0)
	log.SetOutput(l)

	l.Info(msg, nil)

	strout := strings.TrimSpace(buf.String())

	expected := fmt.Sprintf("%s %s \"%s\" %s=%s", level, prefix, msg, key, value)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

}

// Test that a new Logger has the correct defaults
func TestNewLoggerDefaults(t *testing.T) {

	// defaults
	prefix := "SAS:"
	data := F{}
	app_data := F{}
	buf := new(bytes.Buffer)

	l := New(buf, prefix, data, app_data)

	// check Logger defaults match

	if l.name != prefix {
		t.Error("prefix should be '", prefix, "' but is '", l.name, "' instead")
	}
}

// Test that subsequent logger calls use the correct level and doesn't use the previous value
func TestLoggerLevelReset(t *testing.T) {

	// defaults
	prefix := "SAS:"
	data := F{}
	app_data := F{}
	info_msg := "test info"
	debug_msg := "test debug"
	buf := new(bytes.Buffer)

	l := New(buf, prefix, data, app_data)
	log.SetFlags(0)
	log.SetOutput(l)

	// Send an INFO log entry
	l.Info(info_msg, nil)
	strout := strings.TrimSpace(buf.String())

	expected := fmt.Sprintf("INFO %s \"%s\"", prefix, info_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

	// Send a DEBUG log entry
	l.Debug(debug_msg, nil)
	strout = strings.TrimSpace(buf.String())
	expected = fmt.Sprintf("%s %s \"%s\"", "DEBUG", prefix, debug_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

	// Send a second INFO log entry
	l.Info(info_msg, nil)
	strout = strings.TrimSpace(buf.String())
	expected = fmt.Sprintf("INFO %s \"%s\"", prefix, info_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}
}
