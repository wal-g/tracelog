package tracelog

import (
	"io"
	"os"
	"time"
)

var WalgTextFormat = "%s: %s %v"
var WalgTextFormatForDebug = "%s: %s %+v"
var WalgTextFormatFields = []string{"level", "time", "message"}
var WalgDefaultWriter = NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields)

func GetFieldValuesForWalg(loggerType LoggerType) func() Fields {
	return func() Fields {
		now := time.Now()
		fields := Fields{
			"level": loggerType,
			"time":  now.Format("2006/01/02 03:04:05.000000"),
		}

		return fields
	}
}

func GetWalgWriter(out io.Writer) LoggerWriter {
	return NewTextWriter(out, WalgTextFormat, WalgTextFormatFields)
}
