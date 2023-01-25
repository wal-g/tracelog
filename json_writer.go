package tracelog

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type jsonWriter struct {
	lock sync.Mutex
	out  io.Writer
}

func (jsonWriter *jsonWriter) Log(fields Fields) {
	jsonWriter.lock.Lock()
	defer jsonWriter.lock.Unlock()
	data, err := json.MarshalIndent(fields, "", "    ")
	if err == nil {
		_, _ = fmt.Fprint(jsonWriter.out, string(data))
	}
}

func NewJsonWriter(out io.Writer) LoggerWriter {
	return &jsonWriter{
		out: out,
	}
}
