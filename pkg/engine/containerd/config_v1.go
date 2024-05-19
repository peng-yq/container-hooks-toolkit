package containerd

import (
	"fmt"

	"container-hooks-toolkit/pkg/engine"

	"github.com/pelletier/go-toml"
)

// ConfigV1 represents a version 1 containerd config
type ConfigV1 Config

var _ engine.Interface = (*ConfigV1)(nil)

// AddRuntime adds a runtime to the containerd config
func (c *ConfigV1) AddRuntime(name string, path string, setAsDefault bool) error {
	if c == nil || c.Tree == nil {
		return fmt.Errorf("config is nil")
	}

	config := *c.Tree

	config.Set("version", int64(1))

	switch runc := config.GetPath([]string{"plugins", "cri", "containerd", "runtimes", "runc"}).(type) {
	case *toml.Tree:
		runc, _ = toml.Load(runc.String())
		config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name}, runc)
	}

	if config.GetPath([]string{"plugins", "cri", "containerd", "runtimes", name}) == nil {
		config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "runtime_type"}, c.RuntimeType)
		config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "runtime_root"}, "")
		config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "runtime_engine"}, "")
		config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "privileged_without_host_devices"}, false)
	}

	if len(c.ContainerAnnotations) > 0 {
		annotations, err := (*Config)(c).getRuntimeAnnotations([]string{"plugins", "cri", "containerd", "runtimes", name, "container_annotations"})
		if err != nil {
			return err
		}
		annotations = append(c.ContainerAnnotations, annotations...)
		config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "container_annotations"}, annotations)
	}

	config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "options", "BinaryName"}, path)
	config.SetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "options", "Runtime"}, path)

	if setAsDefault && c.UseDefaultRuntimeName {
		config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime_name"}, name)
	} else if setAsDefault {
		// Note: This is deprecated in containerd 1.4.0 and will be removed in 1.5.0
		if config.GetPath([]string{"plugins", "cri", "containerd", "default_runtime"}) == nil {
			config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime", "runtime_type"}, c.RuntimeType)
			config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime", "runtime_root"}, "")
			config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime", "runtime_engine"}, "")
			config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime", "privileged_without_host_devices"}, false)
		}
		config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime", "options", "BinaryName"}, path)
		config.SetPath([]string{"plugins", "cri", "containerd", "default_runtime", "options", "Runtime"}, path)
	}

	*c.Tree = config
	return nil
}

// DefaultRuntime returns the default runtime for the containerd config
func (c ConfigV1) DefaultRuntime() string {
	if runtime, ok := c.GetPath([]string{"plugins", "cri", "containerd", "default_runtime_name"}).(string); ok {
		return runtime
	}
	return ""
}

// RemoveRuntime removes a runtime from the containerd config
func (c *ConfigV1) RemoveRuntime(name string) error {
	if c == nil || c.Tree == nil {
		return nil
	}

	config := *c.Tree

	// If the specified runtime was set as the default runtime we need to remove the default runtime too.
	runtimePath, ok := config.GetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "options", "BinaryName"}).(string)
	if !ok || runtimePath == "" {
		runtimePath, _ = config.GetPath([]string{"plugins", "cri", "containerd", "runtimes", name, "options", "Runtime"}).(string)
	}
	defaultRuntimePath, ok := config.GetPath([]string{"plugins", "cri", "containerd", "default_runtime", "options", "BinaryName"}).(string)
	if !ok || defaultRuntimePath == "" {
		defaultRuntimePath, _ = config.GetPath([]string{"plugins", "cri", "containerd", "default_runtime", "options", "Runtime"}).(string)
	}
	if runtimePath != "" && defaultRuntimePath != "" && runtimePath == defaultRuntimePath {
		config.DeletePath([]string{"plugins", "cri", "containerd", "default_runtime"})
	}

	config.DeletePath([]string{"plugins", "cri", "containerd", "runtimes", name})
	if runtime, ok := config.GetPath([]string{"plugins", "cri", "containerd", "default_runtime_name"}).(string); ok {
		if runtime == name {
			config.DeletePath([]string{"plugins", "cri", "containerd", "default_runtime_name"})
		}
	}

	runtimeConfigPath := []string{"plugins", "cri", "containerd", "runtimes", name}
	for i := 0; i < len(runtimeConfigPath); i++ {
		if runtimes, ok := config.GetPath(runtimeConfigPath[:len(runtimeConfigPath)-i]).(*toml.Tree); ok {
			if len(runtimes.Keys()) == 0 {
				config.DeletePath(runtimeConfigPath[:len(runtimeConfigPath)-i])
			}
		}
	}

	if len(config.Keys()) == 1 && config.Keys()[0] == "version" {
		config.Delete("version")
	}

	*c.Tree = config
	return nil
}

// Save wrotes the config to a file
func (c ConfigV1) Save(path string) (int64, error) {
	return (Config)(c).Save(path)
}
