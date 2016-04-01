saslog
======

A custom GO logger.

Usage
=====

Here's an example usage of saslog.

{{{
package main

import (
    "log"

    "github.com/josepht/saslog/saslog"
)

func main() {
    // Create the logger with a prefix of 'SAS:', a default level of 'INFO',
    // and no logger level data.
    logger := saslog.New("SAS:", "INFO", saslog.F{})

    // Set the logger level data which is included in all logs.
    logger.Data["service"] = "SAS"
    logger.Data["unit"] = "us-sas-1"

    // Turn off flags, saslog composes log messages in logfmt
    log.SetFlags(0)
    // Use our logger as for the standard logger as well
    log.SetOutput(logger)

    // Use our logger directly to include extra data in the logs
    logger.Info("Info", saslog.F{
        "file":  "blah",
        "error": "none",
    })

    // Extra data is not required
    logger.Debug("Debug", nil)

    // The standar logger will use our logger and include our logger level data and log level.
    log.Println("Log message here")
}

}}}
