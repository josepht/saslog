package saslog

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

var buf *bytes.Buffer = new(bytes.Buffer)
var config Config = Config{Writer: buf, Prefix: "SAS:", Name: "sas"}

// Test that the standard logger has our prefix and a default level of 'INFO'
func TestLoggerStdLogOutput(t *testing.T) {
	msg := "test"

	l, err := New(config)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	log.SetFlags(0)
	log.SetOutput(l)

	// Use the standard logger
	log.Print(msg)

	strout := strings.TrimSpace(buf.String())

	expected := fmt.Sprintf("INFO %s \"%s\"", config.Prefix, msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

}

// Test that the loggers output is formated properly and includes the default values
func TestLoggerOutput(t *testing.T) {
	// defaults
	key := "request_id"
	value := "1234"
	msg := "test"

	c := config
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

	strout := strings.TrimSpace(buf.String())

	expected := fmt.Sprintf("INFO %s \"%s\" %s=%s", c.Prefix, msg, key, value)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
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

	l, err := New(config)
	if err != nil {
		t.Error("Failed to create a Logger")
	}

	log.SetFlags(0)
	log.SetOutput(l)

	// Send an INFO log entry
	l.Info(info_msg, nil)
	strout := strings.TrimSpace(buf.String())

	expected := fmt.Sprintf("INFO %s \"%s\"", config.Prefix, info_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

	// Send a DEBUG log entry
	l.Debug(debug_msg, nil)
	strout = strings.TrimSpace(buf.String())
	expected = fmt.Sprintf("%s %s \"%s\"", "DEBUG", config.Prefix, debug_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

	// Send a second INFO log entry
	l.Info(info_msg, nil)
	strout = strings.TrimSpace(buf.String())
	expected = fmt.Sprintf("INFO %s \"%s\"", config.Prefix, info_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}
}
