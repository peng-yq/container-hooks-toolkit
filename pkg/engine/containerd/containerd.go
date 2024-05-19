package containerd

import (
	"container-hooks-toolkit/internal/logger"
	"container-hooks-toolkit/pkg/engine"

	"github.com/pelletier/go-toml"
)

// Config represents the containerd config
type Config struct {
	*toml.Tree
	RuntimeType           string
	UseDefaultRuntimeName bool
	ContainerAnnotations  []string
}

// New creates a containerd config with the specified options
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
