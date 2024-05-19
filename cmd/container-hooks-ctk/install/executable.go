package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type executable struct {
	source string
	target string
}

// install installs e.source to destFolder
func (e executable) install(destFolder string) (string, error) {
	installedfileName, err := installFileToFolderWithName(destFolder, e.target, e.source)
	if err != nil {
		return "", err
	}

	return installedfileName, nil
}

// installFileToFolderWithNam implements "cp src destFolder/name"
func installFileToFolderWithName(destFolder string, name, src string) (string, error) {
	dest := filepath.Join(destFolder, name)
	err := installFile(dest, src)
	if err != nil {
		return "", err
	}
	return dest, nil
}

// installFile copies a file from src to dest and maintains file modes
func installFile(dest string, src string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	err = applyModeFromSource(dest, src)
	if err != nil {
		return err
	}
	return nil
}

// applyModeFromSource sets the file mode for a destination file
// to match that of a specified source file
func applyModeFromSource(dest string, src string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}
	return nil
}

// delete removes the container hooks toolkit from /usr/bin
func delete() error {
	files := []string{
		"/usr/bin/container-hooks-runtime",
		"/usr/bin/container-hooks-ctk",
	}
	for _, file := range files {
		_, err := os.Stat(file)
		if err == nil {
			err = os.Remove(file)
			if err != nil {
				return fmt.Errorf("deleting %s failed (%s)", filepath.Base(file), err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("checking %s failed (%s)", filepath.Base(file), err)
		}
	}
	return nil
}
