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

func TestLogger_PrintWithTextWriter(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("%s [%d] %s: %s", now, os.Getpid(), logLevel, message)
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewTextWriter(&buffer, TestBasicFormat, TestBasicFields))
	infoPostgresLogger.Print(message)
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be: %s, got: %s", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintlnWithTextWriter(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("%s [%d] %s: %s%s", now, os.Getpid(), logLevel, message, "\n")
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewTextWriter(&buffer, TestBasicFormat, TestBasicFields))
	infoPostgresLogger.Println(message)
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be: %v, got: %v", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintfWithTextWriter(t *testing.T) {
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("%s [%d] %s: %s %s", now, os.Getpid(), logLevel, "kek", "test")
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewTextWriter(&buffer, TestBasicFormat, TestBasicFields))
	infoPostgresLogger.Printf("%s %s", "kek", "test")
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be: %v, got: %v", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintWithJsonWriter(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("{\n    \"level\": \"%s\",\n    \"message\": \"%s\",\n    \"pid\": %d,\n    \"time\": \"%s\"\n}", logLevel, message, os.Getpid(), now)
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewJsonWriter(&buffer))
	infoPostgresLogger.Print(message)
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be:\n %s, \ngot: \n%s", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintlnWithJsonWriter(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("{\n    \"level\": \"%s\",\n    \"message\": \"%s\\n\",\n    \"pid\": %d,\n    \"time\": \"%s\"\n}", logLevel, message, os.Getpid(), now)
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewJsonWriter(&buffer))
	infoPostgresLogger.Println(message)
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be:\n %s, \ngot: \n%s", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintfWithJsonWriter(t *testing.T) {
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("{\n    \"level\": \"%s\",\n    \"message\": \"%s\",\n    \"pid\": %d,\n    \"time\": \"%s\"\n}", logLevel, "kek test", os.Getpid(), now)
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewJsonWriter(&buffer))
	infoPostgresLogger.Printf("%s %s", "kek", "test")
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be:\n %s, \ngot: \n%s", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintWithCsvWriter(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("%s,%d,%s,%s\n", now, os.Getpid(), logLevel, message)
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewCsvWriter(&buffer, TestBasicFields))
	infoPostgresLogger.Print(message)
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be:\n %s \ngot:\n %s", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintlnWithCsvWriter(t *testing.T) {
	message := "test"
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("%s,%d,%s,\"%s\n\"\n", now, os.Getpid(), logLevel, message)
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewCsvWriter(&buffer, TestBasicFields))
	infoPostgresLogger.Println(message)
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be:\n %s \ngot:\n %s", expectedOutput, buffer.String())
	}
}

func TestLogger_PrintfWithCsvWriter(t *testing.T) {
	logLevel := InfoLoggerType
	now := time.Now().UTC().Format("2006-01-02 03:04:05 UTC")
	expectedOutput := fmt.Sprintf("%s,%d,%s,%s\n", now, os.Getpid(), logLevel, "kek test")
	buffer := bytes.Buffer{}
	infoPostgresLogger := NewPostgresLogger(GetFieldValuesForTest(logLevel), NewCsvWriter(&buffer, TestBasicFields))
	infoPostgresLogger.Printf("%s %s", "kek", "test")
	if expectedOutput != buffer.String() {
		t.Errorf("\nOutput should be:\n %s \ngot:\n %s", expectedOutput, buffer.String())
	}
}
