package unconfigure

import (
	"fmt"

	"container-hooks-toolkit/internal/logger"
	"container-hooks-toolkit/pkg/engine"
	"container-hooks-toolkit/pkg/engine/containerd"
	"container-hooks-toolkit/pkg/engine/crio"
	"container-hooks-toolkit/pkg/engine/docker"

	cli "github.com/urfave/cli/v2"
)

const (
	defaultRuntime = "docker"

	// defaultRuntimeName is the default name to use in configs for the container hooks runtime
	defaultRuntimeName = "container-hooks-runtime"

	defaultContainerdConfigFilePath = "/etc/containerd/config.toml"
	defaultCrioConfigFilePath       = "/etc/crio/crio.conf"
	defaultDockerConfigFilePath     = "/etc/docker/daemon.json"
)

type command struct {
	logger logger.Interface
}

// NewCommand constructs an configure command with the specified logger
func NewCommand(logger logger.Interface) *cli.Command {
	c := command{
		logger: logger,
	}
	return c.build()
}

// config defines the options that can be set for the CLI through config files,
// environment variables, or command line config
type config struct {
	runtime        string
	configFilePath string

	trustedRuntime struct {
		name string
	}
}

func (m command) build() *cli.Command {
	// Create a config struct to hold the parsed environment variables or command line flags
	config := config{}

	// Create the 'unconfigure' command
	configure := cli.Command{
		Name:  "unconfigure",
		Usage: "Delete trusted container runtime to the specified container engine",
		Before: func(c *cli.Context) error {
			return m.validateFlags(c, &config)
		},
		Action: func(c *cli.Context) error {
			return m.unconfigureConfigFile(c, &config)
		},
	}

	configure.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "runtime",
			Usage:       "Target runtime engine includes [docker, containerd, crio]",
			Value:       defaultRuntime,
			Destination: &config.runtime,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Usage:       "Path to the config file for the target runtime\n\tUse default config path of the target runtime engine if not specified",
			Destination: &config.configFilePath,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "Specify the name of the trusted container runtime that already be added",
			Value:       defaultRuntimeName,
			Destination: &config.trustedRuntime.name,
		},
	}

	return &configure
}

func (m command) validateFlags(c *cli.Context, config *config) error {
	switch config.runtime {
	case "containerd", "crio", "docker":
		break
	default:
		return fmt.Errorf("unrecognized runtime '%v'", config.runtime)
	}
	return nil
}

// unconfigureConfigFile updates the specified container engine config to delete the trusted container runtime
func (m command) unconfigureConfigFile(c *cli.Context, config *config) error {
	configFilePath := config.resolveConfigFilePath()

	var cfg engine.Interface
	var err error
	switch config.runtime {
	case "containerd":
		cfg, err = containerd.New(
			containerd.WithLogger(m.logger),
			containerd.WithPath(configFilePath),
		)
	case "crio":
		cfg, err = crio.New(
			crio.WithLogger(m.logger),
			crio.WithPath(configFilePath),
		)
	case "docker":
		cfg, err = docker.New(
			docker.WithLogger(m.logger),
			docker.WithPath(configFilePath),
		)
	default:
		err = fmt.Errorf("unrecognized runtime '%v'", config.runtime)
	}
	if err != nil || cfg == nil {
		return fmt.Errorf("unable to load config for runtime %v: %v", config.runtime, err)
	}

	err = cfg.RemoveRuntime(config.trustedRuntime.name)
	if err != nil {
		return fmt.Errorf("unable to update config to delete trusted container runtime: %v", err)
	}

	outputPath := config.getOuputConfigPath()
	n, err := cfg.Save(outputPath)
	if err != nil {
		return fmt.Errorf("unable to flush config: %v", err)
	}

	if outputPath != "" {
		if n == 0 {
			m.logger.Infof("Removed empty config from %v", outputPath)
		} else {
			m.logger.Infof("Wrote updated config to %v", outputPath)
		}
		m.logger.Infof("It is recommended that %v daemon be restarted.", config.runtime)
	}

	return nil
}

// resolveConfigFilePath returns the default config file path for the configured container engine
func (c *config) resolveConfigFilePath() string {
	if c.configFilePath != "" {
		return c.configFilePath
	}
	switch c.runtime {
	case "containerd":
		return defaultContainerdConfigFilePath
	case "crio":
		return defaultCrioConfigFilePath
	case "docker":
		return defaultDockerConfigFilePath
	}
	return ""
}

// getOuputConfigPath returns the configured config path
func (c *config) getOuputConfigPath() string {
	return c.resolveConfigFilePath()
}
