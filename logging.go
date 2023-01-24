package tracelog

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

const (
	NormalLogLevel = "NORMAL"
	DevelLogLevel  = "DEVEL"
	ErrorLogLevel  = "ERROR"
)

const (
	CsvPg    = "CSV PG"
	JsonPg   = "JSON PG"
	TextPg   = "TEXT PG"
	TextWalg = "TEXT WALG"
)

var InfoLogger = NewLogger(GetFieldValuesForWalg(InfoLoggerType), WalgDefaultWriter)
var WarningLogger = NewLogger(GetFieldValuesForWalg(WarningLoggerType), WalgDefaultWriter)
var ErrorLogger = NewLogger(GetFieldValuesForWalg(ErrorLoggerType), WalgDefaultWriter)
var DebugLogger = NewLogger(GetFieldValuesForWalg(DebugLoggerType), NewTextWriter(ioutil.Discard, WalgTextFormatForDebug, WalgTextFormatFields))

var LogLevels = []string{NormalLogLevel, DevelLogLevel, ErrorLogLevel}
var logLevel = NormalLogLevel
var logLevelFormatters = map[string]string{
	NormalLogLevel: "%v",
	ErrorLogLevel:  "%v",
	DevelLogLevel:  "%+v",
}

var LogWriters = []string{CsvPg, JsonPg, TextPg, TextWalg}
var logWriter = TextWalg

func setupLoggers() {
	if logWriter == TextWalg {
		setupWalgLoggers()
	} else if logWriter == JsonPg {
		setupJsonPgLoggers()
	} else if logWriter == CsvPg {
		setupCsvPgLoggers()
	} else if logWriter == TextPg {
		setupTextPgLoggers()
	}
}

type LogLevelError struct {
	error
}

func NewLogLevelError(incorrectLogLevel string) LogLevelError {
	return LogLevelError{errors.Errorf("got incorrect log level: '%s', expected one of: '%v'", incorrectLogLevel, LogLevels)}
}

type LogWriterError struct {
	error
}

func NewLogWriterError(incorrectLogWriter string) LogWriterError {
	return LogWriterError{errors.Errorf("got incorrect log writer: '%s', expected one of: '%v'", incorrectLogWriter, strings.Join(LogWriters[:], ","))}
}

func (err LogLevelError) Error() string {
	return fmt.Sprintf(GetErrorFormatter(), err.error)
}

func GetErrorFormatter() string {
	return logLevelFormatters[logLevel]
}

func UpdateLogLevel(newLevel string) error {
	isCorrect := false
	for _, level := range LogLevels {
		if newLevel == level {
			isCorrect = true
		}
	}
	if !isCorrect {
		return NewLogLevelError(newLevel)
	}

	logLevel = newLevel
	setupLoggers()
	return nil
}

func UpdateLogWriter(newWriter string) error {
	isCorrect := false
	for _, logWriter := range LogWriters {
		if newWriter == logWriter {
			isCorrect = true
		}
	}
	if !isCorrect {
		return NewLogWriterError(newWriter)
	}

	logWriter = newWriter
	setupLoggers()
	return nil
}
