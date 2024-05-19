package config

import (
	"container-hooks-toolkit/internal/config"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetOutput returns the effective output
func GetOutput() string {
	return config.GetConfigFilePath()
}

// CreateFolder wraps os.MkdirAll
func CreateFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

// EnsureOutputFolder creates the output folder if it does not exist.
// If the output folder is not specified (i.e. output to STDOUT), it is ignored.
func EnsureOutputFolder() error {
	output := GetOutput()
	if output == "" {
		return nil
	}
	if dir := filepath.Dir(output); dir != "" {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// CreateOutput creates the writer for the output.
func CreateOutput() (io.WriteCloser, error) {
	output := GetOutput()
	if output == "" {
		return nullCloser{os.Stdout}, nil
	}

	return os.Create(output)
}

// note: nullCloser and Close is for config show in stdout (if need)
// nullCloser is a writer that does nothing on Close, actually it's STDOUT
type nullCloser struct {
	io.Writer
}

// Close is a no-op for a nullCloser.
func (d nullCloser) Close() error {
	return nil
}

type configToml config.Toml

var errInvalidConfigOption = errors.New("invalid config option")
var errInvalidFormat = errors.New("invalid format")

// setFlagToKeyValue converts a --set flag to a key-value pair.
// The set flag is of the form key[=value], with the value being optional if key refers to a
// boolean config option.
func (c *configToml) setFlagToKeyValue(setFlag string) (string, interface{}, error) {
	if c == nil {
		return "", nil, errInvalidConfigOption
	}

	setParts := strings.SplitN(setFlag, "=", 2)
	key := setParts[0]

	v := (*config.Toml)(c).Get(key)
	if v == nil {
		return key, nil, errInvalidConfigOption
	}
	switch v.(type) {
	case bool:
		if len(setParts) == 1 {
			return key, true, nil
		}
	}

	if len(setParts) != 2 {
		return key, nil, fmt.Errorf("%w: expected key=value; got %v", errInvalidFormat, setFlag)
	}

	value := setParts[1]
	switch vt := v.(type) {
	case bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return key, value, fmt.Errorf("%w: %w", errInvalidFormat, err)
		}
		return key, b, err
	case string:
		return key, value, nil
	case []string:
		return key, strings.Split(value, ","), nil
	default:
		return key, nil, fmt.Errorf("unsupported type for %v (%v)", setParts, vt)
	}
}
