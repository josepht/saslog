package saslog

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Inspired by http://stackoverflow.com/a/36140590 Thu Mar 31 12:01:08 EDT 2016
type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprint(os.Stderr, time.Now().UTC().Format("2006-01-02 15:04:05.000")+string(bytes))
}

type Logger struct {
	l *log.Logger
}

func (logger *Logger) Init() {
	logger.l = log.New(new(logWriter), " SAS: ", 0)
}

func (logger Logger) log(level string, msg string, data map[string]string) {
	output := "\"" + msg + "\""
	for key, value := range data {
		output += " " + key + "=" + value
	}
	logger.l.Print(level + " " + output)
}

func (logger Logger) Info(msg string, data map[string]string) {
	logger.log("INFO", msg, data)
}
func (logger Logger) Debug(msg string, data map[string]string) {
	logger.log("DEBUG", msg, data)
}

func main() {
	logger := new(Logger)
	logger.Init()

	logger.Info("Info", map[string]string{
		"file":  "blah",
		"error": "none",
	})
	logger.Debug("Debug", nil)
}
