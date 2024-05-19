package oci

import (
	"fmt"
	"os"

	"container-hooks-toolkit/internal/logger"
)

// pathRuntime wraps the path that a binary and defines the semanitcs for how to exec into it.
// This can be used to wrap an OCI-compliant low-level runtime binary, allowing it to be used through the
// Runtime internface.
type pathRuntime struct {
	logger      logger.Interface
	path        string
	execRuntime Runtime
}

var _ Runtime = (*pathRuntime)(nil)

// NewRuntimeForPath creates a Runtime for the specified logger and path
func NewRuntimeForPath(logger logger.Interface, path string) (Runtime, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path '%v': %v", path, err)
	}
	if info.IsDir() || info.Mode()&0111 == 0 {
		return nil, fmt.Errorf("specified path '%v' is not an executable file", path)
	}

	shim := pathRuntime{
		logger:      logger,
		path:        path,
		execRuntime: syscallExec{},
	}

	return &shim, nil
}

// Exec exces into the binary at the path from the pathRuntime struct, passing it the supplied arguments
// after ensuring that the first argument is the path of the target binary.
func (s pathRuntime) Exec(args []string) error {
	runtimeArgs := []string{s.path}
	if len(args) > 1 {
		runtimeArgs = append(runtimeArgs, args[1:]...)
	}

	return s.execRuntime.Exec(runtimeArgs)
}
