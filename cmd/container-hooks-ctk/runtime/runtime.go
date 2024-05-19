package runtime

import (
	"container-hooks-toolkit/cmd/container-hooks-ctk/runtime/configure"
	"container-hooks-toolkit/cmd/container-hooks-ctk/runtime/unconfigure"
	"container-hooks-toolkit/internal/logger"

	"github.com/urfave/cli/v2"
)

type runtimeCommand struct {
	logger logger.Interface
}

// NewCommand constructs a runtime command with the specified logger
func NewCommand(logger logger.Interface) *cli.Command {
	c := runtimeCommand{
		logger: logger,
	}
	return c.build()
}

func (m runtimeCommand) build() *cli.Command {
	// Create the 'runtime' command
	runtime := cli.Command{
		Name:  "runtime",
		Usage: "Configure (Unconfigure) the container hooks runtime",
	}

	runtime.Subcommands = []*cli.Command{
		configure.NewCommand(m.logger),
		unconfigure.NewCommand(m.logger),
	}

	return &runtime
}
