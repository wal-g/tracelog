package tracelog

import (
	"io"
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

func GetPgTextWriter(out io.Writer) LoggerWriter {
	return NewTextWriter(out, PgFormat, PgFields)
}

func GetPgCsvWriter(out io.Writer) LoggerWriter {
	return NewCsvWriter(out, PgFields)
}
