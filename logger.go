package tracelog

import (
	"fmt"
	"os"
)

type LoggerType string

const (
	InfoLoggerType    LoggerType = "INFO"
	WarningLoggerType LoggerType = "WARNING"
	ErrorLoggerType   LoggerType = "ERROR"
	DebugLoggerType   LoggerType = "DEBUG"
)

type Fields map[string]interface{}
type FieldValues func() Fields

type Logger struct {
	loggerWriters []LoggerWriter
	fieldValues   FieldValues
}

func NewLogger(fieldValues FieldValues, loggerWriters ...LoggerWriter) *Logger {
	return &Logger{
		loggerWriters: loggerWriters,
		fieldValues:   fieldValues,
	}
}

func (logger *Logger) Log(v ...interface{}) {
	fields := logger.fieldValues()
	fields["message"] = fmt.Sprint(v...)
	for _, loggerWriter := range logger.loggerWriters {
		loggerWriter.Log(fields)
	}
}

func (logger *Logger) Logf(format string, v ...interface{}) {
	logger.Log(fmt.Sprintf(format, v...))
}

func (logger *Logger) Fatalln(v ...interface{}) {
	logger.Log(fmt.Sprintln(v...))
	os.Exit(1)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.Logf(format, v...)
	os.Exit(1)
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.Log(v...)
	os.Exit(1)
}

func (logger *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	logger.Log(s)
	panic(s)
}

func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Log(s)
	panic(s)
}

func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logger.Log(s)
	panic(s)
}

func (logger *Logger) Println(v ...interface{}) {
	logger.Log(fmt.Sprintln(v...))
}

func (logger *Logger) Printf(format string, v ...interface{}) {
	logger.Logf(format, v...)
}

func (logger *Logger) Print(v ...interface{}) {
	logger.Log(v...)
}
