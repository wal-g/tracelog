package tracelog

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	NormalLogLevel = "NORMAL"
	DevelLogLevel  = "DEVEL"
	ErrorLogLevel  = "ERROR"

	// LegacyLogFormat is the classic "LEVEL: date time message" line format.
	LegacyLogFormat = "LEGACY"
	// TextLogFormat emits key=value records via slog.TextHandler.
	TextLogFormat = "TEXT"
	// JSONLogFormat emits one JSON object per record via slog.JSONHandler.
	JSONLogFormat = "JSON"

	timeFlags = log.LstdFlags | log.Lmicroseconds
)

var InfoLogger = NewErrorLogger(os.Stderr, "INFO: ")
var WarningLogger = NewErrorLogger(os.Stderr, "WARNING: ")
var ErrorLogger = NewErrorLogger(os.Stderr, "ERROR: ")
var DebugLogger = NewErrorLogger(io.Discard, "DEBUG: ")

var LogLevels = []string{NormalLogLevel, DevelLogLevel, ErrorLogLevel}
var LogFormats = []string{LegacyLogFormat, TextLogFormat, JSONLogFormat}

var logLevel = NormalLogLevel
var logFormat = LegacyLogFormat
var logLevelFormatters = map[string]string{
	NormalLogLevel: "%v",
	ErrorLogLevel:  "%v",
	DevelLogLevel:  "%+v",
}

var slogLevels = map[string]slog.Level{
	NormalLogLevel: slog.LevelInfo,
	DevelLogLevel:  slog.LevelDebug,
	ErrorLogLevel:  slog.LevelError,
}

var legacyPrefixes = map[slog.Level]string{
	slog.LevelDebug: "DEBUG: ",
	slog.LevelInfo:  "INFO: ",
	slog.LevelWarn:  "WARNING: ",
	slog.LevelError: "ERROR: ",
}

// slogLevel is shared by every handler created by this package, so the
// minimum level is consistent between tracelog loggers and slog.Default().
var slogLevel = new(slog.LevelVar)

// Standard library logger settings captured right before a TEXT/JSON setup
// redirects it, used to undo that redirection when switching back to LEGACY.
// stdLogBridge is the writer slog.SetDefault installed (nil when not
// redirected); if log.Writer() no longer returns it, the application has
// reconfigured the standard logger since and its settings take precedence.
var (
	stdLogWriter io.Writer
	stdLogFlags  int
	stdLogPrefix string
	stdLogBridge io.Writer
)

// Setup configures the loggers to write to file at the given level, keeping
// the current log format (LEGACY unless SetupWithFormat was called).
func Setup(file *os.File, newLevel string) error {
	return SetupWithFormat(file, newLevel, logFormat)
}

// SetupWithFormat configures the loggers to write to file at the given level
// using one of LogFormats. It also installs a matching slog.Default(), so
// structured attributes can be logged with slog directly. On error the
// current configuration is left unchanged.
func SetupWithFormat(file *os.File, newLevel string, newFormat string) error {
	format, ok := parseLogFormat(newFormat)
	if !ok {
		return NewLogFormatError(newFormat)
	}
	level, ok := slogLevels[newLevel]
	if !ok {
		return newLogLevelError(newLevel)
	}

	logLevel = newLevel
	logFormat = format
	slogLevel.Set(level)

	handler := newHandler(file, slogLevel)
	loggerFor := func(l slog.Level) *errorLogger {
		if l < level {
			return newLogger(io.Discard, l)
		}
		return newHandlerLogger(handler, file, l)
	}
	DebugLogger = loggerFor(slog.LevelDebug)
	InfoLogger = loggerFor(slog.LevelInfo)
	WarningLogger = loggerFor(slog.LevelWarn)
	ErrorLogger = loggerFor(slog.LevelError)

	setSlogDefault(handler)
	return nil
}

func parseLogFormat(format string) (string, bool) {
	format = strings.ToUpper(strings.TrimSpace(format))
	if format == "" {
		return LegacyLogFormat, true
	}
	for _, f := range LogFormats {
		if f == format {
			return f, true
		}
	}
	return "", false
}

// newHandler returns the slog.Handler for the current log format writing to out
// records at or above level.
func newHandler(out io.Writer, level slog.Leveler) slog.Handler {
	switch logFormat {
	case JSONLogFormat:
		return slog.NewJSONHandler(out, handlerOptions(level))
	case TextLogFormat:
		return slog.NewTextHandler(out, handlerOptions(level))
	default:
		return newLegacyHandler(out, level)
	}
}

func handlerOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:       level,
		AddSource:   logLevel == DevelLogLevel,
		ReplaceAttr: replaceTime,
	}
}

// replaceTime renders the record time as UTC ISO 8601 with milliseconds.
func replaceTime(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 && a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
		return slog.String(slog.TimeKey, a.Value.Time().UTC().Format("2006-01-02T15:04:05.000Z"))
	}
	return a
}

// newLogger returns a logger for level l writing to out in the current format.
// Like the LEGACY loggers, it is not gated by the configured log level, so an
// output set explicitly via Set*Output is always written.
func newLogger(out io.Writer, l slog.Level) *errorLogger {
	if out == io.Discard {
		return NewErrorLogger(io.Discard, legacyPrefixes[l])
	}
	return newHandlerLogger(newHandler(out, l), out, l)
}

// newHandlerLogger returns a logger for level l that sends records to handler,
// or writes legacy lines directly to out when the format is LEGACY.
func newHandlerLogger(handler slog.Handler, out io.Writer, l slog.Level) *errorLogger {
	if logFormat == LegacyLogFormat {
		// Bypass the handler so legacy output stays byte-for-byte identical,
		// including SetPrefix/SetFlags behaviour of the embedded log.Logger.
		return NewErrorLogger(out, legacyPrefixes[l])
	}
	w := &slogWriter{handler: handler, level: l, addSource: logLevel == DevelLogLevel}
	return &errorLogger{log.New(w, "", 0)}
}

// setSlogDefault installs handler as slog's default logger. For TEXT and JSON
// this also routes the standard log package through handler; for LEGACY the
// standard log package keeps (or gets back) its original settings.
func setSlogDefault(handler slog.Handler) {
	if logFormat != LegacyLogFormat {
		if stdLogBridge == nil || log.Writer() != stdLogBridge {
			stdLogWriter, stdLogFlags, stdLogPrefix = log.Writer(), log.Flags(), log.Prefix()
		}
		slog.SetDefault(slog.New(handler))
		stdLogBridge = log.Writer()
		return
	}

	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	if stdLogBridge != nil && writer == stdLogBridge {
		writer, flags, prefix = stdLogWriter, stdLogFlags, stdLogPrefix
	}
	stdLogBridge = nil
	slog.SetDefault(slog.New(handler))
	log.SetOutput(writer)
	log.SetFlags(flags)
	log.SetPrefix(prefix)
}

type LogLevelError struct {
	error
}

func NewLogLevelError() LogLevelError {
	return newLogLevelError(logLevel)
}

func newLogLevelError(level string) LogLevelError {
	return LogLevelError{errors.Errorf("got incorrect log level: '%s', expected one of: '%v'", level, LogLevels)}
}

func (err LogLevelError) Error() string {
	return fmt.Sprintf(GetErrorFormatter(), err.error)
}

type LogFormatError struct {
	error
}

func NewLogFormatError(format string) LogFormatError {
	return LogFormatError{errors.Errorf("got incorrect log format: '%s', expected one of: '%v'", format, LogFormats)}
}

func (err LogFormatError) Error() string {
	return fmt.Sprintf(GetErrorFormatter(), err.error)
}

func GetErrorFormatter() string {
	return logLevelFormatters[logLevel]
}

func SetInfoOutput(destination *os.File) {
	InfoLogger = newLogger(destination, slog.LevelInfo)
}

func SetWarningOutput(destination *os.File) {
	WarningLogger = newLogger(destination, slog.LevelWarn)
}

func SetErrorOutput(destination *os.File) {
	ErrorLogger = newLogger(destination, slog.LevelError)
}
