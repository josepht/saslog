package saslog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Test that the standard logger has our prefix and a default level of 'INFO'
func TestLoggerStdLogOutput(t *testing.T) {
	prefix := "SAS:"
	level := "INFO"
	data := F{}
	msg := "test"

	l := New(prefix, level, data)
	log.SetFlags(0)
	log.SetOutput(l)

	old := os.Stderr

	temp, err := ioutil.TempFile(os.TempDir(), "stderr")
	fname := temp.Name()
	if err != nil {
		t.Error("failed to create file '"+fname+"': ", err)
		return
	}

	os.Stderr = temp

	defer func() { os.Stderr = old }()
	defer func() { os.Remove(fname) }()

	// Use the standard logger
	log.Print(msg)
	temp.Close()
	os.Stderr = old

	out, err := ioutil.ReadFile(fname)
	strout := strings.TrimSpace(string(out))

	if err != nil {
		t.Error("failed to read file '"+fname+"': ", err)
		return
	}

	expected := fmt.Sprintf("%s %s \"%s\"", level, prefix, msg)

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
	msg := "test"

	l := New(prefix, level, data)
	log.SetFlags(0)
	log.SetOutput(l)

	temp, err := ioutil.TempFile(os.TempDir(), "stderr")
	fname := temp.Name()
	if err != nil {
		t.Error("failed to create file '"+fname+"': ", err)
		return
	}
	defer func() { os.Remove(fname) }()

	old := os.Stderr

	os.Stderr = temp
	defer func() { os.Stderr = old }()

	l.Info(msg, nil)
	temp.Close()
	os.Stderr = old

	out, err := ioutil.ReadFile(fname)
	strout := strings.TrimSpace(string(out))

	if err != nil {
		t.Error("failed to read file '"+fname+"': ", err)
		return
	}

	expected := fmt.Sprintf("%s %s \"%s\" %s=%s", level, prefix, msg, key, value)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

}

// Test that a new Logger has the correct defaults
func TestNewLoggerDefaults(t *testing.T) {

	// defaults
	prefix := "SAS:"
	level := "INFO"
	data := F{}

	l := New(prefix, level, data)

	// check Logger defaults match
	if l.level != "INFO" {
		t.Error("level should be '", level, "' but is '", l.level, "' instead")
	}

	if l.prefix != prefix {
		t.Error("prefix should be '", prefix, "' but is '", l.prefix, "' instead")
	}

	eq := reflect.DeepEqual(l.Data, data)
	if !eq {
		t.Error("data should be '", data, "' but is '", l.Data, "' instead")
	}
}

// Test that subsequent logger calls use the correct level and doesn't use the previous value
func TestLoggerLevelReset(t *testing.T) {

	// defaults
	prefix := "SAS:"
	level := "INFO"
	data := F{}
	info_msg := "test info"
	debug_msg := "test debug"

	l := New(prefix, level, data)
	log.SetFlags(0)
	log.SetOutput(l)

	// Send an INFO log entry
	temp, err := ioutil.TempFile(os.TempDir(), "stderr")
	fname := temp.Name()
	if err != nil {
		t.Error("failed to create file '"+fname+"': ", err)
		return
	}
	defer func() { os.Remove(fname) }()

	old := os.Stderr

	os.Stderr = temp
	defer func() { os.Stderr = old }()

	l.Info(info_msg, nil)
	temp.Close()
	os.Stderr = old

	out, err := ioutil.ReadFile(fname)
	strout := strings.TrimSpace(string(out))

	if err != nil {
		t.Error("failed to read file '"+fname+"': ", err)
		return
	}

	expected := fmt.Sprintf("%s %s \"%s\"", level, prefix, info_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

	// Send a DEBUG log entry
	temp, err = ioutil.TempFile(os.TempDir(), "stderr")
	fname = temp.Name()
	if err != nil {
		t.Error("failed to create file '"+fname+"': ", err)
		return
	}
	defer func() { os.Remove(fname) }()

	old = os.Stderr

	os.Stderr = temp
	defer func() { os.Stderr = old }()

	l.Debug(debug_msg, nil)
	temp.Close()
	os.Stderr = old

	out, err = ioutil.ReadFile(fname)
	strout = strings.TrimSpace(string(out))

	if err != nil {
		t.Error("failed to read file '"+fname+"': ", err)
		return
	}

	expected = fmt.Sprintf("%s %s \"%s\"", "DEBUG", prefix, debug_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}

	// Send a second INFO log entry
	temp, err = ioutil.TempFile(os.TempDir(), "stderr")
	fname = temp.Name()
	if err != nil {
		t.Error("failed to create file '"+fname+"': ", err)
		return
	}
	defer func() { os.Remove(fname) }()

	old = os.Stderr

	os.Stderr = temp
	defer func() { os.Stderr = old }()

	l.Info(info_msg, nil)
	temp.Close()
	os.Stderr = old

	out, err = ioutil.ReadFile(fname)
	strout = strings.TrimSpace(string(out))

	if err != nil {
		t.Error("failed to read file '"+fname+"': ", err)
		return
	}

	expected = fmt.Sprintf("%s %s \"%s\"", level, prefix, info_msg)

	if !strings.Contains(strout, expected) {
		t.Error(strout, " doesn't containe '"+expected+"'")
	}
}
