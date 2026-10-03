# tracelog

Leveled loggers (`InfoLogger`, `WarningLogger`, `ErrorLogger`, `DebugLogger`) used by WAL-G.

## Configuration

```go
// Level: NORMAL, DEVEL or ERROR. Format: LEGACY (default), TEXT or JSON.
err := tracelog.SetupWithFormat(os.Stderr, tracelog.NormalLogLevel, tracelog.JSONLogFormat)
```

`Setup(file, level)` keeps the current format.

| Format   | Example                                                                      |
|----------|------------------------------------------------------------------------------|
| `LEGACY` | `INFO: 2025/01/02 15:04:05.123456 Backup created`                            |
| `TEXT`   | `time=2025-01-02T15:04:05.123Z level=INFO msg="Backup created"`              |
| `JSON`   | `{"time":"2025-01-02T15:04:05.123Z","level":"INFO","msg":"Backup created"}`  |

In `TEXT` and `JSON` the timestamp is UTC, and at `DEVEL` level every record carries a `source`
attribute pointing at the calling code.

The loggers keep the `*log.Logger` API (`Printf`, `Fatalf`, `SetOutput`, ...), so existing call
sites work unchanged. Setup also installs a matching `slog.Default()`, which lets new code attach
structured attributes:

```go
slog.Info("Backup created", "name", backupName, "size", size)
```

In `TEXT`/`JSON`, output of the standard `log` package goes through the same handler. In `LEGACY`
the standard `log` package is not modified.

### Notes for `TEXT`/`JSON`

- Each logger writes through a bridge into the slog handler. `XLogger.SetOutput(w)` replaces that bridge, so
  anything other than `io.Discard` or a writer saved earlier from `XLogger.Writer()` gets bare messages with no time
  or level.
- Each `Write` call on `XLogger.Writer()` becomes one record. A multi-line chunk stays in one record, and a partial
  line becomes its own record.
- The standard `log` package logs at `INFO`, so at `ERROR` level its output is dropped.
