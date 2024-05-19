package oci

import (
	"fmt"
	"os"
	"syscall"
)

type syscallExec struct{}

var _ Runtime = (*syscallExec)(nil)

func (r syscallExec) Exec(args []string) error {
	err := syscall.Exec(args[0], args, os.Environ())
	if err != nil {
		return fmt.Errorf("could not exec '%v': %v", args[0], err)
	}

	// syscall.Exec is not expected to return. This is an error state regardless of whether
	// err is nil or not.
	return fmt.Errorf("unexpected return from exec '%v'", args[0])
}
