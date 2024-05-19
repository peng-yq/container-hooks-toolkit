package runtime

import (
	"fmt"

	"container-hooks-toolkit/internal/config"
	"container-hooks-toolkit/internal/logger"
	"container-hooks-toolkit/internal/modifier"
	"container-hooks-toolkit/internal/oci"
)

// newContainerRuntime is a factory method that constructs a runtime based on the selected configuration and specified logger
func newContainerRuntime(logger logger.Interface, cfg *config.Config, argv []string) (oci.Runtime, error) {
	lowLevelRuntime, err := oci.NewLowLevelRuntime(logger, cfg.ContainerHookRuntimeConfig.Runtimes)
	if err != nil {
		return nil, fmt.Errorf("error constructing low-level runtime: %v", err)
	}

	if !oci.HasCreateSubcommand(argv) && !oci.HasStartSubcommand(argv) {
		logger.Debugf("Skipping modifier for non-create/start subcommand")
		return lowLevelRuntime, nil
	}

	ociSpec, err := oci.NewSpec(logger, argv)
	if err != nil {
		return nil, fmt.Errorf("error constructing OCI specification: %v", err)
	}

	var specModifier oci.SpecModifier
	if oci.HasCreateSubcommand(argv) {
		logger.Debugf("container runtime mode: create")
		specModifier, err = newSpecModifier(logger, "create", cfg.ContainerRuntimeHookConfig.Path)
	} else if oci.HasStartSubcommand(argv) {
		logger.Debugf("container runtime mode: start")
		specModifier, err = newSpecModifier(logger, "start", cfg.ContainerRuntimeHookConfig.Path)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to construct OCI spec modifier: %v", err)
	}

	// Create the wrapping runtime with the specified modifier
	r := oci.NewModifyingRuntimeWrapper(
		logger,
		lowLevelRuntime,
		ociSpec,
		specModifier,
	)

	return r, nil
}

// newSpecModifier is a factory method that creates constructs an OCI spec modifer based on the provided config.
func newSpecModifier(logger logger.Interface, mode string, hookPath string) (oci.SpecModifier, error) {
	modifier, err := newModifier(logger, mode, hookPath)
	if err != nil {
		return nil, err
	}
	return modifier, nil
}

func newModifier(logger logger.Interface, mode string, hookPath string) (oci.SpecModifier, error) {
	return modifier.NewStableRuntimeModifier(logger, mode, hookPath), nil
}
