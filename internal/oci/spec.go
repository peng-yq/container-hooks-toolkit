package oci

import (
	"fmt"

	"container-hooks-toolkit/internal/logger"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// SpecModifier defines an interace for modifying a (raw) OCI spec
type SpecModifier interface {
	// Modify is a method that accepts a pointer to an OCI Spec and returns an
	// error. The intention is that the function would modify the spec in-place.
	Modify(*specs.Spec) error
}

// Spec defines the operations to be performed on an OCI specification
type Spec interface {
	Load() (*specs.Spec, error)
	Flush() error
	Modify(SpecModifier) error
	LookupEnv(string) (string, bool)
}

// NewSpec creates fileSpec based on the command line arguments passed to the
// application using the specified logger.
func NewSpec(logger logger.Interface, args []string) (Spec, error) {
	bundleDir, err := GetBundleDir(args)
	if err != nil {
		return nil, fmt.Errorf("error getting bundle directory: %v", err)
	}
	logger.Debugf("Using bundle directory: %v", bundleDir)

	ociSpecPath := GetSpecFilePath(bundleDir)
	logger.Infof("Using OCI specification file path: %v", ociSpecPath)

	ociSpec := NewFileSpec(ociSpecPath)

	return ociSpec, nil
}
