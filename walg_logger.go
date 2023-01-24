package tracelog

import (
	"io/ioutil"
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

func setupWalgLoggers() {
	if logLevel == NormalLogLevel {
		DebugLogger = NewLogger(GetFieldValuesForWalg(DebugLoggerType), NewTextWriter(ioutil.Discard, WalgTextFormat, WalgTextFormatFields))
		InfoLogger = NewLogger(GetFieldValuesForWalg(InfoLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
		WarningLogger = NewLogger(GetFieldValuesForWalg(WarningLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
		ErrorLogger = NewLogger(GetFieldValuesForWalg(ErrorLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
	} else if logLevel == ErrorLogLevel {
		ErrorLogger = NewLogger(GetFieldValuesForWalg(ErrorLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
		DebugLogger = NewLogger(GetFieldValuesForWalg(DebugLoggerType), NewTextWriter(ioutil.Discard, WalgTextFormat, WalgTextFormatFields))
		InfoLogger = NewLogger(GetFieldValuesForWalg(InfoLoggerType), NewTextWriter(ioutil.Discard, WalgTextFormat, WalgTextFormatFields))
		WarningLogger = NewLogger(GetFieldValuesForWalg(WarningLoggerType), NewTextWriter(ioutil.Discard, WalgTextFormat, WalgTextFormatFields))
	} else {
		DebugLogger = NewLogger(GetFieldValuesForWalg(DebugLoggerType), NewTextWriter(os.Stdout, WalgTextFormat, WalgTextFormatFields))
		InfoLogger = NewLogger(GetFieldValuesForWalg(InfoLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
		WarningLogger = NewLogger(GetFieldValuesForWalg(WarningLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
		ErrorLogger = NewLogger(GetFieldValuesForWalg(ErrorLoggerType), NewTextWriter(os.Stderr, WalgTextFormat, WalgTextFormatFields))
	}
}
