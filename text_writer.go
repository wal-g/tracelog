package tracelog

import (
	"fmt"
	"io"
	"sync"
)

type textWriter struct {
	lock   sync.Mutex
	out    io.Writer
	format string
	fields []string
}

func (textWriter *textWriter) Log(fields Fields) {
	vals := getValuesFromMap(fields, textWriter.fields)
	textWriter.lock.Lock()
	defer textWriter.lock.Unlock()
	_, _ = fmt.Fprintf(textWriter.out, textWriter.format, vals...)
}

func NewTextWriter(out io.Writer, format string, fields []string) LoggerWriter {
	return &textWriter{
		out:    out,
		format: format,
		fields: fields,
	}
}
