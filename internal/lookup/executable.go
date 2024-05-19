package lookup

import (
	"fmt"
	"os"
	"strings"

	"container-hooks-toolkit/internal/logger"
)

type executable struct {
	file
}

// NewExecutableLocator creates a locator to fine executable files in the path. A logger can also be specified.
func NewExecutableLocator(logger logger.Interface, root string) Locator {
	paths := GetPaths(root)

	return newExecutableLocator(logger, root, paths...)
}

func newExecutableLocator(logger logger.Interface, root string, paths ...string) *executable {
	f := newFileLocator(
		WithLogger(logger),
		WithRoot(root),
		WithSearchPaths(paths...),
		WithFilter(assertExecutable),
		WithCount(1),
	)

	l := executable{
		file: *f,
	}

	return &l
}

var _ Locator = (*executable)(nil)

// Locate finds executable files with the specified pattern in the path.
// If a relative or absolute path is specified, the prefix paths are not considered.
func (p executable) Locate(pattern string) ([]string, error) {
	// For absolute paths we ensure that it is executable
	if strings.Contains(pattern, "/") {
		err := assertExecutable(pattern)
		if err != nil {
			return nil, fmt.Errorf("absolute path %v is not an executable file: %v", pattern, err)
		}
		return []string{pattern}, nil
	}
	// For relative path, not ensure it is executable
	return p.file.Locate(pattern)
}

// assertExecutable checks whether the specified path is an execuable file.
func assertExecutable(filename string) error {
	err := assertFile(filename)
	if err != nil {
		return err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if info.Mode()&0111 == 0 {
		return fmt.Errorf("specified file '%v' is not executable", filename)
	}

	return nil
}
