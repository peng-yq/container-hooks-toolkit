package oci

// Runtime is an interface for a runtime shim. The Exec method accepts a list
// of command line arguments, and returns an error / nil.
type Runtime interface {
	Exec([]string) error
}
