package crio

import (
	"fmt"
	"os"

	"container-hooks-toolkit/internal/logger"

	"github.com/pelletier/go-toml"
)

type builder struct {
	logger logger.Interface
	path   string
}

// Option defines a function that can be used to configure the config builder
type Option func(*builder)

// WithLogger sets the logger for the config builder
func WithLogger(logger logger.Interface) Option {
	return func(b *builder) {
		b.logger = logger
	}
}

// WithPath sets the path for the config builder
func WithPath(path string) Option {
	return func(b *builder) {
		b.path = path
	}
}

func (b *builder) build() (*Config, error) {
	if b.path == "" {
		empty := toml.Tree{}
		return (*Config)(&empty), nil
	}
	if b.logger == nil {
		b.logger = logger.New()
	}

	return b.loadConfig(b.path)
}

// loadConfig loads the cri-o config from disk
func (b *builder) loadConfig(config string) (*Config, error) {
	b.logger.Infof("Loading config: %v", config)

	info, err := os.Stat(config)
	if os.IsExist(err) && info.IsDir() {
		return nil, fmt.Errorf("config file is a directory")
	}

	if os.IsNotExist(err) {
		b.logger.Infof("Config file does not exist; using empty config")
		config = "/dev/null"
	} else {
		b.logger.Infof("Loading config from %v", config)
	}

	cfg, err := toml.LoadFile(config)
	if err != nil {
		return nil, err
	}

	b.logger.Infof("Successfully loaded config")

	return (*Config)(cfg), nil
}
