package gcppubsubdemo

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// LogType represents the type of logging to perform.
// Current types are DisableLogLeve, ErrorLogLevel,
// DebugLogLevel, and TraceLogLevel.
type LogType int

const (
	// DisableLogLevel represents no logging
	DisableLogLevel LogType = iota
	// ErrorLogLevel represents error log entries
	ErrorLogLevel
	// DebugLogLevel represents debug log entries
	DebugLogLevel
	// TraceLogLevel represents trace log entries
	TraceLogLevel
)

var logger struct {
	TRACE *log.Logger
	DEBUG *log.Logger
	ERROR *log.Logger

	loggers [TraceLogLevel]*log.Logger
	output  io.Writer
	level   LogType
}

// LogLevel sets the log level. The supported levels are
// error (1), debug (2), trace (3). A level of 0 disables
// all logging. A level of 3 or higher enables everything.
func LogLevel(level LogType) {
	logger.level = level
	for i := range logger.loggers {
		if i+1 <= int(level) {
			logger.loggers[i].SetOutput(logger.output)
		} else {
			logger.loggers[i].SetOutput(ioutil.Discard)
		}
	}
}

// LogOutput sets the logging output to an io.Writer. This allows
// the API user to direct log output to a file, for example.
func LogOutput(writer io.Writer) {
	logger.output = writer
	LogLevel(logger.level)
}

func init() {
	var prefix string

	for i := range logger.loggers {
		switch LogType(i + 1) {
		case ErrorLogLevel:
			prefix = "ERROR "
		case DebugLogLevel:
			prefix = "DEBUG "
		case TraceLogLevel:
			prefix = "TRACE "
		}

		logger.loggers[i] = log.New(ioutil.Discard, prefix, log.Ldate|log.Ltime)
	}

	logger.TRACE = logger.loggers[TraceLogLevel-1]
	logger.DEBUG = logger.loggers[DebugLogLevel-1]
	logger.ERROR = logger.loggers[ErrorLogLevel-1]

	logger.output = os.Stderr
}
