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

type errorLogger struct {
	internalLogger *log.Logger
	onFatal        func(format string, err error)
	onPanic        func(format string, err error)
}

func NewErrorLogger(out io.Writer, loggerType LoggerType) *errorLogger {
	logPrefix := fmt.Sprintf("%s: ", loggerType)
	internalLogger := log.New(out, logPrefix, timeFlags)
	onFatal := func(format string, err error) {
		os.Exit(1)
	}
	onPanic := func(format string, err error) {
		panic(fmt.Sprintf(format, err))
	}
	return &errorLogger{internalLogger: internalLogger, onFatal: onFatal, onPanic: onPanic}
}

func (logger *errorLogger) SetOutput(w io.Writer) {
	logger.internalLogger.SetOutput(w)
}

func (logger *errorLogger) PanicError(err error) {
	logger.internalLogger.Printf(GetErrorFormatter(), err)
	logger.onPanic(GetErrorFormatter(), err)
}

func (logger *errorLogger) PanicfOnError(format string, err error) {
	if err != nil {
		logger.internalLogger.Printf(format, err)
		logger.onPanic(format, err)
	}
}

func (logger *errorLogger) PanicOnError(err error) {
	if err != nil {
		logger.PanicError(err)
	}
}

func (logger *errorLogger) FatalError(err error) {
	logger.internalLogger.Printf(GetErrorFormatter(), err)
	logger.onFatal(GetErrorFormatter(), err)
}

func (logger *errorLogger) FatalfOnError(format string, err error) {
	if err != nil {
		logger.internalLogger.Printf(format, err)
		logger.onFatal(format, err)
	}
}

func (logger *errorLogger) FatalOnError(err error) {
	if err != nil {
		logger.FatalError(err)
	}
}

func (logger *errorLogger) PrintError(err error) {
	logger.internalLogger.Printf(GetErrorFormatter()+"\n", err)
}

func (logger *errorLogger) PrintOnError(err error) {
	if err != nil {
		logger.PrintError(err)
	}
}

func (logger *errorLogger) Printf(format string, v ...interface{}) {
	logger.internalLogger.Printf(format, v)
}

func (logger *errorLogger) Println(v ...interface{}) {
	logger.internalLogger.Println(v)
}
