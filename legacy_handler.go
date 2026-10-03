package tracelog

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const legacyTimeFormat = "2006/01/02 15:04:05.000000"

// legacyHandler is a slog.Handler that renders records the same way the
// LEGACY tracelog loggers do ("INFO: 2006/01/02 15:04:05.000000 message"),
// followed by any attributes as " key=value" pairs.
type legacyHandler struct {
	w      io.Writer
	level  slog.Leveler
	mu     *sync.Mutex // shared by clones created via WithAttrs/WithGroup
	attrs  []byte      // pre-rendered attributes from WithAttrs
	prefix string      // group prefix from WithGroup, e.g. "db.pool."
}

func newLegacyHandler(w io.Writer, level slog.Leveler) *legacyHandler {
	return &legacyHandler{w: w, level: level, mu: &sync.Mutex{}}
}

func (h *legacyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *legacyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	buf := bytes.NewBuffer(append([]byte(nil), h.attrs...))
	for _, a := range attrs {
		appendLegacyAttr(buf, h.prefix, a)
	}
	h2.attrs = buf.Bytes()
	return &h2
}

func (h *legacyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	h2.prefix = h.prefix + name + "."
	return &h2
}

//nolint:gocritic // slog.Handler requires Record to be passed by value.
func (h *legacyHandler) Handle(_ context.Context, r slog.Record) error {
	buf := &bytes.Buffer{}
	buf.WriteString(legacyLevelPrefix(r.Level))
	buf.WriteString(r.Time.Format(legacyTimeFormat))
	buf.WriteByte(' ')
	buf.WriteString(r.Message)
	buf.Write(h.attrs)
	r.Attrs(func(a slog.Attr) bool {
		appendLegacyAttr(buf, h.prefix, a)
		return true
	})
	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	return err
}

func legacyLevelPrefix(level slog.Level) string {
	if prefix, ok := legacyPrefixes[level]; ok {
		return prefix
	}
	return level.String() + ": "
}

func appendLegacyAttr(buf *bytes.Buffer, prefix string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		if a.Key != "" {
			prefix += a.Key + "."
		}
		for _, ga := range a.Value.Group() {
			appendLegacyAttr(buf, prefix, ga)
		}
		return
	}
	buf.WriteByte(' ')
	buf.WriteString(prefix)
	buf.WriteString(a.Key)
	buf.WriteByte('=')
	buf.WriteString(quoteIfNeeded(a.Value.String()))
}

func quoteIfNeeded(s string) string {
	if s == "" || strings.IndexFunc(s, needsQuote) >= 0 {
		return strconv.Quote(s)
	}
	return s
}

// needsQuote reports whether r would make an unquoted value ambiguous or let it
// corrupt the line: spaces, '"', '=', and control or non-printable characters.
func needsQuote(r rune) bool {
	return r == ' ' || r == '"' || r == '=' || !unicode.IsPrint(r)
}
