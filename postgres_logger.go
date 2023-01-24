package tracelog

import (
	"io/ioutil"
	"os"
	"time"
)

var PgFormat = "%s [%d] %s: %s"
var PgFields = []string{"timestamp", "pid", "error_severity", "message"}

func GetFieldValuesForPg(loggerType LoggerType) func() Fields {
	return func() Fields {
		now := time.Now().UTC()
		fields := Fields{
			"timestamp":      now.Format("2006-01-02 03:04:05.000 UTC"),
			"pid":            os.Getpid(),
			"error_severity": loggerType,
		}

		return fields
	}
}

func setupJsonPgLoggers() {
	if logLevel == NormalLogLevel {
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewJsonWriter(ioutil.Discard))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewJsonWriter(os.Stderr))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewJsonWriter(os.Stderr))
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewJsonWriter(os.Stderr))
	} else if logLevel == ErrorLogLevel {
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewJsonWriter(os.Stderr))
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewJsonWriter(ioutil.Discard))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewJsonWriter(ioutil.Discard))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewJsonWriter(ioutil.Discard))
	} else {
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewJsonWriter(os.Stdout))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewJsonWriter(os.Stderr))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewJsonWriter(os.Stderr))
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewJsonWriter(os.Stderr))
	}
}

func setupTextPgLoggers() {
	if logLevel == NormalLogLevel {
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewTextWriter(ioutil.Discard, PgFormat, PgFields))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
	} else if logLevel == ErrorLogLevel {
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewTextWriter(ioutil.Discard, PgFormat, PgFields))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewTextWriter(ioutil.Discard, PgFormat, PgFields))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewTextWriter(ioutil.Discard, PgFormat, PgFields))
	} else {
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewTextWriter(os.Stdout, PgFormat, PgFields))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewTextWriter(os.Stderr, PgFormat, PgFields))
	}
}

func setupCsvPgLoggers() {
	if logLevel == NormalLogLevel {
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewCsvWriter(ioutil.Discard, PgFields))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewCsvWriter(os.Stderr, PgFields))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewCsvWriter(os.Stderr, PgFields))
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewCsvWriter(os.Stderr, PgFields))
	} else if logLevel == ErrorLogLevel {
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewCsvWriter(os.Stderr, PgFields))
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewCsvWriter(ioutil.Discard, PgFields))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewCsvWriter(ioutil.Discard, PgFields))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewCsvWriter(ioutil.Discard, PgFields))
	} else {
		DebugLogger = NewLogger(GetFieldValuesForPg(DebugLoggerType), NewCsvWriter(os.Stdout, PgFields))
		InfoLogger = NewLogger(GetFieldValuesForPg(InfoLoggerType), NewCsvWriter(os.Stderr, PgFields))
		WarningLogger = NewLogger(GetFieldValuesForPg(WarningLoggerType), NewCsvWriter(os.Stderr, PgFields))
		ErrorLogger = NewLogger(GetFieldValuesForPg(ErrorLoggerType), NewCsvWriter(os.Stderr, PgFields))
	}
}
