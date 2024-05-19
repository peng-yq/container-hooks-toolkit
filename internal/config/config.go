package config

import (
	"path/filepath"
)

const (
	configFilePath 	   = "container-hooks/config.toml"
	CTKExecutable      = "container-hooks-ctk"
	CTKDefaultFilePath = "/usr/bin/container-hooks-ctk"
	ContainerHookPath = "/etc/container-hooks/hooks.json"
)

var (
	// DefaultExecutableDir specifies the default path to use for executables if they cannot be located in the path.
	DefaultExecutableDir = "/usr/bin"
)

// Config represents the contents of the config.toml file for the Trusted Container Toolkit
type Config struct {
	ContainerHookCtkConfig              CTKConfig         `toml:"container-hooks-ctk"`
	ContainerHookRuntimeConfig     	    RuntimeConfig     `toml:"container-hooks-runtime"`
	ContainerRuntimeHookConfig          RuntimeHookConfig `toml:"container-hooks"`
}

// GetConfigFilePath returns the path to the config file for the container-hooks-runtime
// /etc/configFilePath
func GetConfigFilePath() string {
	return filepath.Join("/etc", configFilePath)
}

// GetConfig sets up the config struct. Values are read from a toml file
// or set via the environment.
func GetConfig() (*Config, error) {
	cfg, err := New(
		WithConfigFile(GetConfigFilePath()),
	)
	if err != nil {
		return nil, err
	}

	return cfg.Config()
}

// GetDefault defines the default values for the config
func GetDefault() (*Config, error) {
	c := Config{
		ContainerHookCtkConfig: CTKConfig{
			Path: CTKDefaultFilePath,
		},
		ContainerHookRuntimeConfig: RuntimeConfig{
			DebugFilePath: "/etc/container-hooks/container-hooks-runtime.log",
			LogLevel:      "info",
			Runtimes:      []string{"runc", "docker-runc"},
		},
		ContainerRuntimeHookConfig: RuntimeHookConfig{
			Path: ContainerHookPath,
		},
	}
	return &c, nil
}