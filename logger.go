package tracelog

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LoggerType string

const (
	InfoLoggerType    LoggerType = "INFO"
	WarningLoggerType LoggerType = "WARNING"
	ErrorLoggerType   LoggerType = "ERROR"
	DebugLoggerType   LoggerType = "DEBUG"
)

// BasicLogger represents the logger
// without the panic and fatal methods
type BasicLogger interface {
	Output(calldepth int, s string) error
	Writer() io.Writer
	SetOutput(w io.Writer)
	Flags() int
	SetFlags(flag int)
	Prefix() string
	SetPrefix(prefix string)
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type logger struct {
	BasicLogger
	onFatal func(format string, v ...interface{})
	onPanic func(format string, v ...interface{})
}

func NewLogger(out io.Writer, loggerType LoggerType) *logger {
	logPrefix := fmt.Sprintf("%s: ", loggerType)
	internalLogger := log.New(out, logPrefix, timeFlags)
	onFatal := func(format string, v ...interface{}) {
		os.Exit(1)
	}
	onPanic := func(format string, v ...interface{}) {
		panic(fmt.Sprintf(format, v...))
	}
	return &logger{BasicLogger: internalLogger, onFatal: onFatal, onPanic: onPanic}
}

func (logger *logger) PanicError(err error) {
	logger.Printf(GetErrorFormatter(), err)
	logger.onPanic(GetErrorFormatter(), err)
}

func (logger *logger) PanicfOnError(format string, err error) {
	if err != nil {
		logger.Printf(format, err)
		logger.onPanic(format, err)
	}
}

func (logger *logger) PanicOnError(err error) {
	if err != nil {
		logger.PanicError(err)
	}
}

func (logger *logger) FatalError(err error) {
	logger.Printf(GetErrorFormatter(), err)
	logger.onFatal(GetErrorFormatter(), err)
}

func (logger *logger) FatalfOnError(format string, err error) {
	if err != nil {
		logger.Printf(format, err)
		logger.onFatal(format, err)
	}
}

func (logger *logger) FatalOnError(err error) {
	if err != nil {
		logger.FatalError(err)
	}
}

func (logger *logger) PrintError(err error) {
	logger.Printf(GetErrorFormatter()+"\n", err)
}

func (logger *logger) PrintOnError(err error) {
	if err != nil {
		logger.PrintError(err)
	}
}

func (logger *logger) Fatal(v ...interface{}) {
	logger.Print(v...)
	logger.onFatal(GetErrorFormatter(), v)
}

func (logger *logger) Fatalf(format string, v ...interface{}) {
	logger.Printf(format, v...)
	logger.onFatal(format, v)
}

func (logger *logger) Fatalln(v ...interface{}) {
	logger.Println(v...)
	logger.onFatal(GetErrorFormatter(), v)
}

func (logger *logger) Panic(v ...interface{}) {
	logger.Print(v...)
	logger.onPanic(GetErrorFormatter(), v)
}

func (logger *logger) Panicf(format string, v ...interface{}) {
	logger.Printf(format, v...)
	logger.onPanic(format, v)
}

func (logger *logger) Panicln(v ...interface{}) {
	logger.Println(v...)
	logger.onPanic(GetErrorFormatter(), v)
}
