package tracelog_test

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/wal-g/tracelog"
)

// setup configures tracelog to log into a temp file and returns a function
// that reads everything logged so far.
func setup(t *testing.T, level, format string) func() string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "log-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := tracelog.SetupWithFormat(os.Stderr, tracelog.NormalLogLevel, tracelog.LegacyLogFormat); err != nil {
			t.Error(err)
		}
		_ = file.Close()
	})
	if err := tracelog.SetupWithFormat(file, level, format); err != nil {
		t.Fatal(err)
	}
	return func() string {
		data, err := os.ReadFile(file.Name())
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
}

func jsonRecords(t *testing.T, out string) []map[string]any {
	t.Helper()
	var records []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if line == "" {
			continue
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("line is not JSON: %q: %v", line, err)
		}
		records = append(records, rec)
	}
	return records
}

func TestLegacyFormat(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.LegacyLogFormat)

	tracelog.InfoLogger.Printf("hello %s", "world")
	tracelog.WarningLogger.Println("careful")
	tracelog.ErrorLogger.PrintError(errors.New("boom"))
	tracelog.DebugLogger.Println("hidden")
	slog.Info("structured", "key", "value", "spaced", "a b")
	slog.Debug("hidden")

	stamp := `\d{4}/\d\d/\d\d \d\d:\d\d:\d\d\.\d{6}`
	want := regexp.MustCompile(`^` +
		`INFO: ` + stamp + ` hello world\n` +
		`WARNING: ` + stamp + ` careful\n` +
		`ERROR: ` + stamp + ` boom\n` +
		`INFO: ` + stamp + ` structured key=value spaced="a b"\n$`)
	if out := read(); !want.MatchString(out) {
		t.Errorf("unexpected legacy output:\n%s", out)
	}
}

func TestLegacyLeavesStdLogUntouched(t *testing.T) {
	setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)
	setup(t, tracelog.NormalLogLevel, tracelog.LegacyLogFormat)
	if log.Writer() != os.Stderr || log.Flags() != log.LstdFlags || log.Prefix() != "" {
		t.Errorf("standard logger was not restored: writer=%v flags=%d prefix=%q", log.Writer(), log.Flags(), log.Prefix())
	}
}

func TestLegacyRestoresStdLogSetAfterInit(t *testing.T) {
	custom, err := os.CreateTemp(t.TempDir(), "std-*")
	if err != nil {
		t.Fatal(err)
	}
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	t.Cleanup(func() {
		log.SetOutput(writer)
		log.SetFlags(flags)
		log.SetPrefix(prefix)
		_ = custom.Close()
	})
	// The application reconfigures the standard logger while it is redirected.
	setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)
	log.SetOutput(custom)
	log.SetPrefix("app: ")

	setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)
	setup(t, tracelog.NormalLogLevel, tracelog.LegacyLogFormat)
	if log.Writer() != custom || log.Prefix() != "app: " {
		t.Errorf("standard logger was not restored: writer=%v prefix=%q", log.Writer(), log.Prefix())
	}
}

func TestLegacyQuotesControlChars(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.LegacyLogFormat)

	slog.Info("m", "cr", "a\rb", "esc", "\x1b[31m")

	if out := read(); !strings.HasSuffix(out, ` m cr="a\rb" esc="\x1b[31m"`+"\n") {
		t.Errorf("control characters not quoted: %q", out)
	}
}

func TestJSONFormat(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, "json")

	tracelog.InfoLogger.Printf("hello %s", "world")
	tracelog.WarningLogger.Println("careful")
	tracelog.ErrorLogger.PrintError(errors.New("multi\nline"))
	tracelog.DebugLogger.Println("hidden")
	slog.Info("structured", "key", 42)
	log.Print("from std log")

	records := jsonRecords(t, read())
	want := []struct{ level, msg string }{
		{"INFO", "hello world"},
		{"WARN", "careful"},
		{"ERROR", "multi\nline"},
		{"INFO", "structured"},
		{"INFO", "from std log"},
	}
	if len(records) != len(want) {
		t.Fatalf("got %d records, want %d: %v", len(records), len(want), records)
	}
	for i, w := range want {
		if records[i]["level"] != w.level || records[i]["msg"] != w.msg {
			t.Errorf("record %d = %v, want level=%s msg=%q", i, records[i], w.level, w.msg)
		}
		if _, ok := records[i]["source"]; ok {
			t.Errorf("record %d has source outside DEVEL: %v", i, records[i])
		}
	}
	if records[3]["key"] != float64(42) {
		t.Errorf("structured attribute missing: %v", records[3])
	}
	if !regexp.MustCompile(`^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d\.\d{3}Z$`).MatchString(records[0]["time"].(string)) {
		t.Errorf("unexpected time format: %v", records[0]["time"])
	}
}

func TestJSONSourcePointsAtCaller(t *testing.T) {
	read := setup(t, tracelog.DevelLogLevel, tracelog.JSONLogFormat)

	tracelog.DebugLogger.Printf("debug")
	tracelog.ErrorLogger.PrintOnError(errors.New("wrapped helper"))

	records := jsonRecords(t, read())
	if len(records) != 2 {
		t.Fatalf("got %d records, want 2: %v", len(records), records)
	}
	for _, rec := range records {
		source, _ := rec["source"].(map[string]any)
		if file, _ := source["file"].(string); filepath.Base(file) != "logging_test.go" {
			t.Errorf("source does not point at the caller: %v", rec)
		}
	}
	if records[0]["level"] != "DEBUG" {
		t.Errorf("debug record missing in DEVEL: %v", records[0])
	}
}

func TestTextFormat(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.TextLogFormat)

	tracelog.InfoLogger.Print("first\nsecond")
	tracelog.WarningLogger.Printf("careful")

	lines := strings.Split(strings.TrimSuffix(read(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), lines)
	}
	if !strings.Contains(lines[0], `level=INFO msg="first\nsecond"`) {
		t.Errorf("unexpected line: %s", lines[0])
	}
	if !strings.Contains(lines[1], `level=WARN msg=careful`) {
		t.Errorf("unexpected line: %s", lines[1])
	}
}

func TestErrorLevelOnlyLogsErrors(t *testing.T) {
	read := setup(t, tracelog.ErrorLogLevel, tracelog.JSONLogFormat)

	tracelog.InfoLogger.Println("info")
	tracelog.WarningLogger.Println("warning")
	slog.Warn("slog warning")
	tracelog.ErrorLogger.Println("error")

	records := jsonRecords(t, read())
	if len(records) != 1 || records[0]["msg"] != "error" {
		t.Errorf("unexpected records: %v", records)
	}
}

func TestLoggerMethodsInStructuredFormat(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)

	tracelog.InfoLogger.SetPrefix("PREFETCH " + tracelog.InfoLogger.Prefix())
	tracelog.InfoLogger.Println("prefixed")
	tracelog.WarningLogger.SetOutput(io.Discard)
	tracelog.WarningLogger.Println("discarded")

	records := jsonRecords(t, read())
	if len(records) != 1 || records[0]["msg"] != "PREFETCH prefixed" {
		t.Errorf("unexpected records: %v", records)
	}
}

func TestPanicOnErrorLogsThenPanics(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)

	func() {
		defer func() {
			if recover() == nil {
				t.Error("PanicOnError did not panic")
			}
		}()
		tracelog.ErrorLogger.PanicOnError(errors.New("fatal"))
	}()

	records := jsonRecords(t, read())
	if len(records) != 1 || records[0]["msg"] != "fatal" || records[0]["level"] != "ERROR" {
		t.Errorf("unexpected records: %v", records)
	}
}

func TestSetWarningOutputFollowsFormat(t *testing.T) {
	setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)
	file, err := os.CreateTemp(t.TempDir(), "warn-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()

	tracelog.SetWarningOutput(file)
	tracelog.WarningLogger.Println("redirected")

	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	records := jsonRecords(t, string(data))
	if len(records) != 1 || records[0]["msg"] != "redirected" {
		t.Errorf("unexpected records: %v", records)
	}
}

func TestSetInfoOutputAtErrorLevel(t *testing.T) {
	for _, format := range tracelog.LogFormats {
		t.Run(format, func(t *testing.T) {
			setup(t, tracelog.ErrorLogLevel, format)
			file, err := os.CreateTemp(t.TempDir(), "info-*")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = file.Close() }()

			tracelog.SetInfoOutput(file)
			tracelog.InfoLogger.Println("redirected")

			data, err := os.ReadFile(file.Name())
			if err != nil {
				t.Fatal(err)
			}
			if lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n"); len(lines) != 1 ||
				!strings.Contains(lines[0], "redirected") {
				t.Errorf("unexpected output: %q", data)
			}
		})
	}
}

func TestInvalidSetupKeepsConfiguration(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)

	err := tracelog.SetupWithFormat(os.Stderr, tracelog.NormalLogLevel, "xml")
	var formatErr tracelog.LogFormatError
	if !errors.As(err, &formatErr) || !strings.Contains(err.Error(), "'xml'") {
		t.Errorf("unexpected error for invalid format: %v", err)
	}
	err = tracelog.Setup(os.Stderr, "VERBOSE")
	var levelErr tracelog.LogLevelError
	if !errors.As(err, &levelErr) || !strings.Contains(err.Error(), "'VERBOSE'") {
		t.Errorf("unexpected error for invalid level: %v", err)
	}

	tracelog.InfoLogger.Println("still json")
	if records := jsonRecords(t, read()); len(records) != 1 {
		t.Errorf("configuration changed after invalid setup: %v", records)
	}
}

func TestSetupKeepsFormat(t *testing.T) {
	read := setup(t, tracelog.NormalLogLevel, tracelog.JSONLogFormat)
	file, err := os.CreateTemp(t.TempDir(), "log-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()

	if err := tracelog.Setup(file, tracelog.DevelLogLevel); err != nil {
		t.Fatal(err)
	}
	tracelog.DebugLogger.Println("debug")

	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if records := jsonRecords(t, string(data)); len(records) != 1 || records[0]["level"] != "DEBUG" {
		t.Errorf("unexpected records: %v", records)
	}
	if out := read(); out != "" {
		t.Errorf("old destination still written: %q", out)
	}
}
