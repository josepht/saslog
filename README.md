saslog
======

A custom GO logger.

Usage
=====

Here's an example usage of saslog.

```
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/josepht/saslog"
)

func main() {
	// Create the logger.
	c := saslog.Config{
		Writer: os.Stderr,
		Name:   "sas",
		Prefix: "SAS:",
		SystemData: saslog.F{
			"service": "sas",
			"unit":    "sas-us-1",
		},
		AppData: saslog.F{"revno": "678"},
	}
	logger, err := saslog.New(c)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	logger.Info("Config based log", nil)

	// Turn off flags, saslog composes log messages in logfmt
	log.SetFlags(0)
	// Use our logger as the standard logger as well
	log.SetOutput(logger)

	// Use our logger directly to include extra data in the logs
	logger.Info("Info", saslog.F{
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
