package engine

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents a runtime config
type Config string

// Write writes the specified contents to a config file.
func (c Config) Write(output []byte) (int, error) {
	path := string(c)
	if path == "" {
		n, err := os.Stdout.Write(output)
		if err == nil {
			os.Stdout.WriteString("\n")
		}
		return n, err
	}

	if len(output) == 0 {
		err := os.Remove(path)
		if err != nil {
			return 0, fmt.Errorf("unable to remove empty file: %v", err)
		}
		return 0, nil
	}

	if dir := filepath.Dir(path); dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return 0, fmt.Errorf("unable to create directory %v: %v", dir, err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("unable to open %v for writing: %v", path, err)
	}
	defer f.Close()

	return f.Write(output)
}
