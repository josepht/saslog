slog
======

A custom GO logger.

Usage
=====

Here's an example usage of slog.

```
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/josepht/slog"
)

func main() {
	// Create the logger.
	c := slog.Config{
		Writer: os.Stderr,
		Name:   "SLOG:",
		Prefix: "slog",
		SystemData: slog.F{
			"service": "slog",
			"unit":    "slog-us-1",
		},
		AppData: slog.F{"revno": "678"},
	}
	logger := slog.New(c)

	// Derive a Logger from an existing one.
	l := logger.New(slog.Config{Name: "SUBSYSTEM:"})
	l.Info("Derived Logger", nil)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Turn off flags, slog composes log messages in logfmt
	log.SetFlags(0)

	// Use our logger as the standard logger as well
	log.SetOutput(logger)

	// Use our logger directly to include extra data in the logs
	logger.Info("Info", slog.F{
		"file":  "blah",
		"error": "none",
	})

	// Extra data is not required
	logger.Debug("Debug", nil)

	// The standard logger will use our logger and include our
	// logger level data and log level.
	log.Println("Log message here")
}
```
