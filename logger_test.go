package tracelog

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

var TestBasicFormat = "%s [%d] %s: %s"
var TestBasicFields = []string{"time", "pid", "level", "message"}

func GetFieldValuesForTest(loggerType LoggerType) func() Fields {
	return func() Fields {
		now := time.Now().UTC()
		fields := Fields{
			"time":  now.Format("2006-01-02 03:04:05 UTC"),
			"pid":   os.Getpid(),
			"level": loggerType,
		}

		return fields
	}
}

func TestLogger_PrintWith(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	buffer := bytes.Buffer{}
	tests := []struct {
		name       string
		input      string
		want       string
		infoLogger *Logger
	}{
		{
			"TextWriter",
			message,
			fmt.Sprintf("%s [%d] %s: %s", now, os.Getpid(), logLevel, message),
			NewLogger(GetFieldValuesForTest(logLevel), NewTextWriter(&buffer, TestBasicFormat, TestBasicFields)),
		},
		{
			"JsonWriter",
			message,
			fmt.Sprintf("{\n    \"level\": \"%s\",\n    \"message\": \"%s\",\n    \"pid\": %d,\n    \"time\": \"%s\"\n}", logLevel, message, os.Getpid(), now),
			NewLogger(GetFieldValuesForTest(logLevel), NewJsonWriter(&buffer)),
		},
		{
			"CsvWriter",
			message,
			fmt.Sprintf("%s,%d,%s,%s\n", now, os.Getpid(), logLevel, message),
			NewLogger(GetFieldValuesForTest(logLevel), NewCsvWriter(&buffer, TestBasicFields)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.infoLogger.Print(tt.input)
			if tt.want != buffer.String() {
				t.Errorf("\\nOutput should be:\\n %s, \\ngot: \\n%s\"", tt.want, buffer.String())
			}
			buffer.Reset()
		})
	}
}

func TestLogger_PrintlnWith(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	buffer := bytes.Buffer{}
	tests := []struct {
		name       string
		input      string
		want       string
		infoLogger *Logger
	}{
		{
			"TextWriter",
			message,
			fmt.Sprintf("%s [%d] %s: %s%s", now, os.Getpid(), logLevel, message, "\n"),
			NewLogger(GetFieldValuesForTest(logLevel), NewTextWriter(&buffer, TestBasicFormat, TestBasicFields)),
		},
		{
			"JsonWriter",
			message,
			fmt.Sprintf("{\n    \"level\": \"%s\",\n    \"message\": \"%s\\n\",\n    \"pid\": %d,\n    \"time\": \"%s\"\n}", logLevel, message, os.Getpid(), now),
			NewLogger(GetFieldValuesForTest(logLevel), NewJsonWriter(&buffer)),
		},
		{
			"CsvWriter",
			message,
			fmt.Sprintf("%s,%d,%s,\"%s\n\"\n", now, os.Getpid(), logLevel, message),
			NewLogger(GetFieldValuesForTest(logLevel), NewCsvWriter(&buffer, TestBasicFields)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.infoLogger.Println(tt.input)
			if tt.want != buffer.String() {
				t.Errorf("\\nOutput should be:\\n %s, \\ngot: \\n%s\"", tt.want, buffer.String())
			}
			buffer.Reset()
		})
	}
}

func TestLogger_PrintfWith(t *testing.T) {
	format := "%s %s"
	messages := []string{"kek", "test"}
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	buffer := bytes.Buffer{}
	tests := []struct {
		name       string
		format     string
		input      []string
		want       string
		infoLogger *Logger
	}{
		{
			"TextWriter",
			format,
			messages,
			fmt.Sprintf("%s [%d] %s: %s %s", now, os.Getpid(), logLevel, "kek", "test"),
			NewLogger(GetFieldValuesForTest(logLevel), NewTextWriter(&buffer, TestBasicFormat, TestBasicFields)),
		},
		{
			"JsonWriter",
			format,
			messages,
			fmt.Sprintf("{\n    \"level\": \"%s\",\n    \"message\": \"%s\",\n    \"pid\": %d,\n    \"time\": \"%s\"\n}", logLevel, "kek test", os.Getpid(), now),
			NewLogger(GetFieldValuesForTest(logLevel), NewJsonWriter(&buffer)),
		},
		{
			"CsvWriter",
			format,
			messages,
			fmt.Sprintf("%s,%d,%s,%s\n", now, os.Getpid(), logLevel, "kek test"),
			NewLogger(GetFieldValuesForTest(logLevel), NewCsvWriter(&buffer, TestBasicFields)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.infoLogger.Printf(tt.format, tt.input[0], tt.input[1])
			if tt.want != buffer.String() {
				t.Errorf("\\nOutput should be:\\n %s, \\ngot: \\n%s\"", tt.want, buffer.String())
			}
			buffer.Reset()
		})
	}
}
