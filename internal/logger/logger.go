package logger

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// text formatter
type LogFormatter struct {
}

func (t LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:06")

	if len(entry.Data) > 0 {
		var dataStr string
		for key, value := range entry.Data {
			dataStr += fmt.Sprintf("%s:%v ", key, value)
		}
		fmt.Fprintf(b, "[%s] %s", timestamp, dataStr)
	} else {
		fmt.Fprintf(b, "[%s]", timestamp)
	}

	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		fmt.Fprintf(b, " %s %s %s \n", fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, " msg:%s \n", entry.Message)
	}

	return b.Bytes(), nil
}

func CreateLogFile(filename string) (*os.File, error) {
	if filename == "" || filename == os.DevNull {
		return nil, nil
	}
	if dir := filepath.Dir(filepath.Clean(filename)); dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
