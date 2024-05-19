package modifier

import (
	"os"
	"errors"
	"encoding/json"
	"path/filepath"

	"container-hooks-toolkit/internal/logger"
	"container-hooks-toolkit/internal/oci"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// NewStableRuntimeModifier creates an OCI spec modifier that inserts the container hooks into an OCI
// spec. The specified logger is used to capture log output.
func NewStableRuntimeModifier(logger logger.Interface, mode string, hookPath string) oci.SpecModifier {
	m := stableRuntimeModifier{
		logger:                          logger,
		mode:                            mode,
		hookPath:						 hookPath,
	}

	return &m
}

// stableRuntimeModifier modifies an OCI spec inplace, inserting the container hooks.
// If the hook is already present, no modification is made.
type stableRuntimeModifier struct {
	logger                 			logger.Interface
	// mode if one of [create, start] depends on the runc args
	mode      	string
	hookPath 	string
}

// Modify applies the required modification to the incoming OCI spec, inserting the container-hooks-runtime-hook.
func (m stableRuntimeModifier) Modify(spec *specs.Spec) error {
    // If an container hooks already exists, we don't make any modifications to the spec.
    if spec.Hooks != nil {
        if m.mode == "create" || m.mode == "start" {
            for _, hook := range spec.Hooks.Prestart {
                if hasRuntimeHook(&hook) {
                    m.logger.Infof("Existing hooks (%v) found in OCI spec", hook.Path)
                    return nil
                }
            }
        }
    }

    if spec.Hooks == nil {
        spec.Hooks = &specs.Hooks{}
    }

    // Read hooks JSON file
    hooksData, err := os.ReadFile(m.hookPath)
    if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// if hooks file is not exist, we don't insert hooks
			return nil
		} else {
			m.logger.Errorf("Failed to read hooks JSON file: %v", err)
			return err
		}
    }

    // Unmarshal JSON data into Hooks struct
    var hooks specs.Hooks
    if err := json.Unmarshal(hooksData, &hooks); err != nil {
        m.logger.Errorf("Failed to unmarshal hooks JSON: %v", err)
        return err
    }

    if m.mode == "create" || m.mode == "start" {
		spec.Hooks.Prestart = append(spec.Hooks.Prestart, hooks.Prestart...)
		spec.Hooks.CreateRuntime = append(spec.Hooks.CreateRuntime, hooks.CreateRuntime...)
		spec.Hooks.CreateContainer = append(spec.Hooks.CreateContainer, hooks.CreateContainer...)
		spec.Hooks.StartContainer = append(spec.Hooks.StartContainer, hooks.StartContainer...)
		spec.Hooks.Poststart = append(spec.Hooks.Poststart, hooks.Poststart...)
		spec.Hooks.Poststop = append(spec.Hooks.Poststop, hooks.Poststop...)
		m.logger.Debugf("Add hooks success")
	}

    return nil
}



// hasRuntimeHook checks if the provided hook has inject container hooks
func hasRuntimeHook(hook *specs.Hook) bool {
	bins := map[string]struct{}{
		"container-hooks": {},
	}

	_, exists := bins[filepath.Base(hook.Path)]

	return exists
}