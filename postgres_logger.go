package tracelog

import (
	"io/ioutil"
	"os"
	"time"
)

type PostgresLogger struct {
	*Logger
}

var BasicFormat = "%s [%d] %s: %s"
var BasicFields = []string{"timestamp", "pid", "error_severity", "message"}

func GetFieldValues(loggerType LoggerType) func() Fields {
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

var InfoPostgresLogger = NewPostgresLogger(GetFieldValues(InfoLoggerType), NewTextWriter(os.Stderr, BasicFormat, BasicFields))
var WarningPostgresLogger = NewPostgresLogger(GetFieldValues(WarningLoggerType), NewTextWriter(os.Stderr, BasicFormat, BasicFields))
var ErrorPostgresLogger = NewPostgresLogger(GetFieldValues(ErrorLoggerType), NewTextWriter(os.Stderr, BasicFormat, BasicFields))
var DebugPostgresLogger = NewPostgresLogger(GetFieldValues(DebugLoggerType), NewTextWriter(ioutil.Discard, BasicFormat, BasicFields))

func NewPostgresLogger(fieldValues FieldValues, loggerWriters ...LoggerWriter) *PostgresLogger {
	return &PostgresLogger{
		NewLogger(fieldValues, loggerWriters...),
	}
}

func (logger *PostgresLogger) PanicError(err error) {
	logger.Panic(err)
}

func (logger *PostgresLogger) PanicfOnError(format string, err error) {
	if err != nil {
		logger.Panicf(format, err)
	}
}

func (logger *PostgresLogger) PanicOnError(err error) {
	if err != nil {
		logger.PanicError(err)
	}
}

func (logger *PostgresLogger) FatalError(err error) {
	logger.Fatal(err)
}

func (logger *PostgresLogger) FatalfOnError(format string, err error) {
	if err != nil {
		logger.Fatalf(format, err)
	}
}

func (logger *PostgresLogger) FatalOnError(err error) {
	if err != nil {
		logger.FatalError(err)
	}
}

func (logger *PostgresLogger) PrintError(err error) {
	logger.Println(err)
}

func (logger *PostgresLogger) PrintOnError(err error) {
	if err != nil {
		logger.PrintError(err)
	}
}
