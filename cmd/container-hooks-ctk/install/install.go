package install

import (
	"container-hooks-toolkit/internal/logger"
	"fmt"
	"path/filepath"

	cli "github.com/urfave/cli/v2"
)

type installCommand struct {
	logger logger.Interface
}

// options stores the subcommand options
type options struct {
	toolkitRoot string
}

// NewCommand constructs a install command with the specified logger
func NewCommand(logger logger.Interface) *cli.Command {
	c := installCommand{
		logger: logger,
	}
	return c.build()
}

func (m installCommand) build() *cli.Command {
	opts := options{}

	// Create the 'install' command
	install := cli.Command{
		Name:  "install",
		Usage: "Install the components of the container hooks toolkit",
		Before: func(c *cli.Context) error {
			return m.validateOptions(c, &opts)
		},
		Action: func(ctx *cli.Context) error {
			return m.install(ctx, &opts)
		},
	}

	install.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "toolkit-root",
			Usage:       "The directory where the container hooks toolkit is installed",
			Destination: &opts.toolkitRoot,
		},
	}

	return &install
}

// validateOptions checks whether the specified options are valid
func (m installCommand) validateOptions(c *cli.Context, opts *options) error {
	if opts.toolkitRoot == "" {
		return fmt.Errorf("--toolkit-root must be specified")
	}
	return nil
}

// install installs the components of the trusted container toolkit to /usr/bin and set PATH.
func (m installCommand) install(cli *cli.Context, opts *options) error {
	m.logger.Infof("Installing container hooks toolkit from '%v'", opts.toolkitRoot)
	m.logger.Infof("Removing existing container hooks toolkit installation")
	err := delete()
	if err != nil {
		return fmt.Errorf("error removing container hooks toolkit: %v", err)
	}

	err = m.installToolkit(opts.toolkitRoot, "container-hooks-runtime")
	if err != nil {
		return fmt.Errorf("error installing container-hooks-runtime: %v", err)
	}

	err = m.installToolkit(opts.toolkitRoot, "container-hooks")
	if err != nil {
		return fmt.Errorf("error installing container-hooks: %v", err)
	}

	err = m.installToolkit(opts.toolkitRoot, "container-hooks-ctk")
	if err != nil {
		return fmt.Errorf("error installing container-hooks-ctk: %v", err)
	}

	return nil
}

// installToolkit installs the container hooks executable and add to PATH.
func (m installCommand) installToolkit(toolkitRoot string, tool string) error {
	m.logger.Infof("Installing %v from '%v'", tool, toolkitRoot)
	targetPath := filepath.Join(toolkitRoot, tool)
	e := executable{
		source: targetPath,
		target: tool,
	}
	_, err := e.install("/usr/bin")
	if err != nil {
		return err
	}
	return nil
}
