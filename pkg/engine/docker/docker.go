package docker

import (
	"encoding/json"
	"fmt"

	"container-hooks-toolkit/internal/logger"
	"container-hooks-toolkit/pkg/engine"
)

const (
	defaultDockerRuntime = "runc"
)

// Config defines a docker config file.
type Config map[string]interface{}

// New creates a docker config with the specified options
func New(opts ...Option) (engine.Interface, error) {
	b := &builder{}
	for _, opt := range opts {
		opt(b)
	}

	if b.logger == nil {
		b.logger = logger.New()
	}

	return b.build()
}

// AddRuntime adds a new runtime to the docker config
func (c *Config) AddRuntime(name string, path string, setAsDefault bool) error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	config := *c

	// Read the existing runtimes
	runtimes := make(map[string]interface{})
	if _, exists := config["runtimes"]; exists {
		runtimes = config["runtimes"].(map[string]interface{})
	}

	// Add / update the runtime definitions
	runtimes[name] = map[string]interface{}{
		"path":        path,
		"runtimeArgs": []string{},
	}

	config["runtimes"] = runtimes

	if setAsDefault {
		config["default-runtime"] = name
	}

	*c = config
	return nil
}

// DefaultRuntime returns the default runtime for the docker config
func (c Config) DefaultRuntime() string {
	r, ok := c["default-runtime"].(string)
	if !ok {
		return ""
	}
	return r
}

// RemoveRuntime removes a runtime from the docker config
func (c *Config) RemoveRuntime(name string) error {
	if c == nil {
		return nil
	}
	config := *c

	if _, exists := config["default-runtime"]; exists {
		defaultRuntime := config["default-runtime"].(string)
		if defaultRuntime == name {
			config["default-runtime"] = defaultDockerRuntime
		}
	}

	if _, exists := config["runtimes"]; exists {
		runtimes := config["runtimes"].(map[string]interface{})

		delete(runtimes, name)

		if len(runtimes) == 0 {
			delete(config, "runtimes")
		}
	}

	*c = config

	return nil
}

// Save writes the config to the specified path
func (c Config) Save(path string) (int64, error) {
	output, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return 0, fmt.Errorf("unable to convert to JSON: %v", err)
	}

	n, err := engine.Config(path).Write(output)
	return int64(n), err
}
