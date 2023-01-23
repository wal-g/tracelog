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

var InfoPostgresLogger = NewLogger(GetFieldValues(InfoLoggerType), NewTextWriter(os.Stderr, BasicFormat, BasicFields))
var WarningPostgresLogger = NewLogger(GetFieldValues(WarningLoggerType), NewTextWriter(os.Stderr, BasicFormat, BasicFields))
var ErrorPostgresLogger = NewLogger(GetFieldValues(ErrorLoggerType), NewTextWriter(os.Stderr, BasicFormat, BasicFields))
var DebugPostgresLogger = NewLogger(GetFieldValues(DebugLoggerType), NewTextWriter(ioutil.Discard, BasicFormat, BasicFields))
