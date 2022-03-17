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

type logger struct {
	internalLogger *log.Logger
	onFatal        func(format string, err error)
	onPanic        func(format string, err error)
}

func NewLogger(out io.Writer, loggerType LoggerType) *logger {
	logPrefix := fmt.Sprintf("%s: ", loggerType)
	internalLogger := log.New(out, logPrefix, timeFlags)
	onFatal := func(format string, err error) {
		os.Exit(1)
	}
	onPanic := func(format string, err error) {
		panic(fmt.Sprintf(format, err))
	}
	return &logger{internalLogger: internalLogger, onFatal: onFatal, onPanic: onPanic}
}

func (logger *logger) SetOutput(w io.Writer) {
	logger.internalLogger.SetOutput(w)
}

func (logger *logger) PanicError(err error) {
	logger.internalLogger.Printf(GetErrorFormatter(), err)
	logger.onPanic(GetErrorFormatter(), err)
}

func (logger *logger) PanicfOnError(format string, err error) {
	if err != nil {
		logger.internalLogger.Printf(format, err)
		logger.onPanic(format, err)
	}
}

func (logger *logger) PanicOnError(err error) {
	if err != nil {
		logger.PanicError(err)
	}
}

func (logger *logger) FatalError(err error) {
	logger.internalLogger.Printf(GetErrorFormatter(), err)
	logger.onFatal(GetErrorFormatter(), err)
}

func (logger *logger) FatalfOnError(format string, err error) {
	if err != nil {
		logger.internalLogger.Printf(format, err)
		logger.onFatal(format, err)
	}
}

func (logger *logger) FatalOnError(err error) {
	if err != nil {
		logger.FatalError(err)
	}
}

func (logger *logger) PrintError(err error) {
	logger.internalLogger.Printf(GetErrorFormatter()+"\n", err)
}

func (logger *logger) PrintOnError(err error) {
	if err != nil {
		logger.PrintError(err)
	}
}

func (logger *logger) Print(v ...interface{}) {
	logger.internalLogger.Print(v)
}

func (logger *logger) Printf(format string, v ...interface{}) {
	logger.internalLogger.Printf(format, v)
}

func (logger *logger) Println(v ...interface{}) {
	logger.internalLogger.Println(v)
}
