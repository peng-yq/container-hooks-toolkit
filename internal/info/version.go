package info

import "strings"

// version must be set by go build's -X main.version= option in the Makefile.
var version = ""

// gitCommit will be the hash that the binary was built from
// and will be populated by the Makefile
var gitCommit = ""

// GetVersionParts returns the different version components
func GetVersionParts() []string {
	v := []string{version}

	if gitCommit != "" {
		v = append(v, "commit: "+gitCommit)
	}

	return v
}

// GetVersionString returns the string representation of the version
func GetVersionString(more ...string) string {
	v := append(GetVersionParts(), more...)
	return strings.Join(v, "\n")
}
