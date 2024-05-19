package oci

import (
	"fmt"

	"container-hooks-toolkit/internal/logger"
)

type modifyingRuntimeWrapper struct {
	logger   logger.Interface
	runtime  Runtime
	ociSpec  Spec
	modifier SpecModifier
}

var _ Runtime = (*modifyingRuntimeWrapper)(nil)

// NewModifyingRuntimeWrapper creates a runtime wrapper that applies the specified modifier to the OCI specification
// before invoking the wrapped runtime. If the modifier is nil, the input runtime is returned.
func NewModifyingRuntimeWrapper(logger logger.Interface, runtime Runtime, spec Spec, modifier SpecModifier) Runtime {
	if modifier == nil {
		logger.Infof("Using low-level runtime with no modification")
		return runtime
	}

	rt := modifyingRuntimeWrapper{
		logger:   logger,
		runtime:  runtime,
		ociSpec:  spec,
		modifier: modifier,
	}
	return &rt
}

// Exec checks whether a modification of the OCI specification is required and modifies it accordingly before exec-ing
// into the wrapped runtime.
func (r *modifyingRuntimeWrapper) Exec(args []string) error {
	if HasCreateSubcommand(args) || HasStartSubcommand(args) {
		err := r.modify()
		if err != nil {
			return fmt.Errorf("could not apply required modification to OCI specification: %v", err)
		}
		r.logger.Infof("Applied required modification to OCI specification")
	} else {
		r.logger.Infof("No modification of OCI specification required")
	}

	r.logger.Infof("Forwarding command to runtime")
	return r.runtime.Exec(args)
}

// modify loads, modifies, and flushes the OCI specification using the defined Modifier
func (r *modifyingRuntimeWrapper) modify() error {
	_, err := r.ociSpec.Load()
	if err != nil {
		return fmt.Errorf("error loading OCI specification for modification: %v", err)
	}

	err = r.ociSpec.Modify(r.modifier)
	if err != nil {
		return fmt.Errorf("error modifying OCI spec: %v", err)
	}

	err = r.ociSpec.Flush()
	if err != nil {
		return fmt.Errorf("error writing modified OCI specification: %v", err)
	}
	return nil
}
