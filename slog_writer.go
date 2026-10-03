package tracelog

import (
	"context"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// packagePath is the function name prefix of this package, e.g. "github.com/wal-g/tracelog.".
var packagePath = reflect.TypeFor[slogWriter]().PkgPath() + "."

// slogWriter is the io.Writer behind an errorLogger's embedded log.Logger in
// TEXT and JSON formats. Every log.Logger.Output call performs exactly one
// Write, which becomes one slog.Record at the logger's level.
type slogWriter struct {
	handler   slog.Handler
	level     slog.Level
	addSource bool // whether to look up the caller for the source attribute
}

func (w *slogWriter) Write(p []byte) (int, error) {
	ctx := context.Background()
	if !w.handler.Enabled(ctx, w.level) {
		return len(p), nil
	}
	msg := p
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	var pc uintptr
	if w.addSource {
		pc = callerPC()
	}
	r := slog.NewRecord(time.Now(), w.level, string(msg), pc)
	return len(p), w.handler.Handle(ctx, r)
}

// callerPC returns the program counter of the first caller outside of the log
// and tracelog packages, so the source attribute points at the code that
// logged rather than at helpers such as errorLogger.FatalOnError.
func callerPC() uintptr {
	var pcs [32]uintptr
	// Skip runtime.Callers, callerPC and slogWriter.Write.
	n := runtime.Callers(3, pcs[:])
	// runtime.Callers returns one PC per logical frame, inlined frames
	// included, so each PC can be resolved on its own.
	for _, pc := range pcs[:n] {
		frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
		if !isLoggingFrame(frame.Function) {
			return pc
		}
	}
	return 0
}

func isLoggingFrame(function string) bool {
	return strings.HasPrefix(function, "log.") || strings.HasPrefix(function, packagePath)
}
