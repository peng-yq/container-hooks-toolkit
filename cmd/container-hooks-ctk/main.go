package main

import (
	"os"

	"container-hooks-toolkit/cmd/container-hooks-ctk/config"
	"container-hooks-toolkit/cmd/container-hooks-ctk/install"
	"container-hooks-toolkit/cmd/container-hooks-ctk/runtime"
	"container-hooks-toolkit/internal/info"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

// options defines the options that can be set for the CLI through config files,
// environment variables, or command line flags
type options struct {
	// Debug indicates whether the CLI is started in "debug" mode
	Debug bool
	// Quiet indicates whether the CLI is started in "quiet" mode
	Quiet bool
}

func main() {
	logger := log.New()

	// Create a options struct to hold the parsed environment variables or command line flags
	opts := options{}

	// Create the top-level CLI
	c := cli.NewApp()
	c.Name = "container-hooks-ctk"
	c.UseShortOptionHandling = true
	c.EnableBashCompletion = true
	c.Usage = "Tools to configure the container hooks toolkit"
	c.Version = info.GetVersionString()

	// Setup the flags for this command
	c.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Enable debug-level logging",
			Destination: &opts.Debug,
		},
		&cli.BoolFlag{
			Name:        "quiet",
			Aliases:     []string{"q"},
			Usage:       "Suppress all output except for errors",
			Destination: &opts.Quiet,
		},
	}

	// Set log-level for all subcommands
	c.Before = func(c *cli.Context) error {
		logLevel := log.InfoLevel
		if opts.Debug {
			logLevel = log.DebugLevel
		}
		if opts.Quiet {
			logLevel = log.ErrorLevel
		}
		logger.SetLevel(logLevel)
		return nil
	}

	// Define the subcommands
	c.Commands = []*cli.Command{
		runtime.NewCommand(logger),
		config.NewCommand(logger),
		install.NewCommand(logger),
	}

	// Run the CLI
	err := c.Run(os.Args)
	if err != nil {
		logger.Errorf("%v", err)
		os.Exit(1)
	}
}
