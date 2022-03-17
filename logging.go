package tracelog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

type LogLevel string

const (
	NormalLogLevel LogLevel = "NORMAL"
	DevelLogLevel  LogLevel = "DEVEL"
	ErrorLogLevel  LogLevel = "ERROR"
	timeFlags               = log.LstdFlags | log.Lmicroseconds
)

var InfoLogger = NewErrorLogger(os.Stdout, InfoLoggerType)
var WarningLogger = NewErrorLogger(os.Stdout, WarningLoggerType)
var ErrorLogger = NewErrorLogger(os.Stderr, ErrorLoggerType)
var DebugLogger = NewErrorLogger(ioutil.Discard, DebugLoggerType)

var LogLevels = []LogLevel{NormalLogLevel, DevelLogLevel, ErrorLogLevel}
var logLevel = NormalLogLevel
var logLevelFormatters = map[LogLevel]string{
	NormalLogLevel: "%v",
	ErrorLogLevel:  "%v",
	DevelLogLevel:  "%+v",
}

func syncLogLevel() {
	switch logLevel {
	case NormalLogLevel:
		DebugLogger.SetOutput(ioutil.Discard)
	case ErrorLogLevel:
		DebugLogger.SetOutput(ioutil.Discard)
		InfoLogger.SetOutput(ioutil.Discard)
		WarningLogger.SetOutput(ioutil.Discard)
	default: // assume DevelLogLevel
		DebugLogger.SetOutput(os.Stdout)
	}
}

type LogLevelError struct {
	error
}

func NewLogLevelError() LogLevelError {
	return LogLevelError{errors.Errorf("got incorrect log level: '%s', expected one of: '%v'", logLevel, LogLevels)}
}

func (err LogLevelError) Error() string {
	return fmt.Sprintf(GetErrorFormatter(), err.error)
}

func GetErrorFormatter() string {
	return logLevelFormatters[logLevel]
}

func UpdateLogLevel(newLevel LogLevel) error {
	isCorrect := false
	for _, level := range LogLevels {
		if newLevel == level {
			isCorrect = true
		}
	}
	if !isCorrect {
		return NewLogLevelError()
	}

	logLevel = newLevel
	syncLogLevel()
	return nil
}

func SetOnPanicFunc(onPanic func(format string, err error)) {
	updateLoggers(func(logger *errorLogger) {
		logger.onPanic = onPanic
	})
}

func SetOnFatalFunc(onFatal func(format string, err error)) {
	updateLoggers(func(logger *errorLogger) {
		logger.onFatal = onFatal
	})
}

func updateLoggers(updateFunc func(*errorLogger)) {
	for _, logger := range []*errorLogger{InfoLogger, WarningLogger, DebugLogger, ErrorLogger} {
		updateFunc(logger)
	}
}
