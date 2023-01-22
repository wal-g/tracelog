package tracelog

import (
	"encoding/csv"
	"fmt"
	"io"
	"sync"
)

type csvWriter struct {
	lock   sync.Mutex
	out    io.Writer
	fields []string
}

func (csvWriter *csvWriter) Log(fields Fields) {
	vals := getValuesFromMap(fields, csvWriter.fields)
	stringVals := make([]string, len(vals))
	for i := range vals {
		stringVals[i] = fmt.Sprint(vals[i])
	}
	csvWriter.lock.Lock()
	defer csvWriter.lock.Unlock()
	w := csv.NewWriter(csvWriter.out)
	err := w.Write(stringVals)
	if err == nil {
		w.Flush()
	}
}

func NewCsvWriter(out io.Writer, fields []string) LoggerWriter {
	return &csvWriter{
		out:    out,
		fields: fields,
	}
}
