package config

import (
	"fmt"

	"container-hooks-toolkit/internal/config"
	"container-hooks-toolkit/internal/logger"

	cli "github.com/urfave/cli/v2"
)

type command struct {
	logger logger.Interface
}

// options stores the subcommand options
type options struct {
	sets cli.StringSlice
}

// NewCommand constructs an config command with the specified logger
func NewCommand(logger logger.Interface) *cli.Command {
	c := command{
		logger: logger,
	}
	return c.build()
}

// build
func (m command) build() *cli.Command {
	opts := options{}

	// Create the 'config' command
	c := cli.Command{
		Name:  "config",
		Usage: "Generate the container hooks toolkit configuration",
		Action: func(ctx *cli.Context) error {
			return run(ctx, &opts)
		},
	}

	c.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "set",
			Usage:       "Set a config value using the pattern key=value.This flag can be specified multiple times",
			Destination: &opts.sets,
		},
	}

	return &c
}

func run(c *cli.Context, opts *options) error {
	// create config.toml
	cfgToml, err := config.New(
		config.WithConfigFile(GetOutput()),
	)
	if err != nil {
		return fmt.Errorf("unable to create config: %v", err)
	}

	for _, set := range opts.sets.Value() {
		key, value, err := (*configToml)(cfgToml).setFlagToKeyValue(set)
		if err != nil {
			return fmt.Errorf("invalid --set option %v: %w", set, err)
		}
		cfgToml.Set(key, value)
	}

	if err := EnsureOutputFolder(); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	output, err := CreateOutput()
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer output.Close()
	cfgToml.Save(output)

	return nil
}
