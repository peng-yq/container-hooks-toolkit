package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

// Toml is a type for the TOML representation of a config.
type Toml toml.Tree

type options struct {
	configFile string
}

// Option is a functional option for loading TOML config files.
type Option func(*options)

// WithConfigFile sets the config file option.
func WithConfigFile(configFile string) Option {
	return func(o *options) {
		o.configFile = configFile
	}
}

// New creates a new toml tree based on the provided options
func New(opts ...Option) (*Toml, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return loadConfigToml(o.configFile)
}

func loadConfigToml(filename string) (*Toml, error) {
	if filename == "" {
		return defaultToml()
	}

	tomlFile, err := os.Open(filename)
	if os.IsNotExist(err) {
		return defaultToml()
	} else if err != nil {
		return nil, fmt.Errorf("failed to load specified config file: %v", err)
	}
	defer tomlFile.Close()

	return loadConfigTomlFrom(tomlFile)

}

func defaultToml() (*Toml, error) {
	cfg, err := GetDefault()
	if err != nil {
		return nil, err
	}
	contents, err := toml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	return loadConfigTomlFrom(bytes.NewReader(contents))
}

func loadConfigTomlFrom(reader io.Reader) (*Toml, error) {
	tree, err := toml.LoadReader(reader)
	if err != nil {
		return nil, err
	}
	return (*Toml)(tree), nil
}

// Config returns the typed config associated with the toml tree.
func (t *Toml) Config() (*Config, error) {
	cfg, err := GetDefault()
	if err != nil {
		return nil, err
	}
	if t == nil {
		return cfg, nil
	}
	if err := t.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return cfg, nil
}

// Unmarshal wraps the toml.Tree Unmarshal function.
func (t *Toml) Unmarshal(v interface{}) error {
	return (*toml.Tree)(t).Unmarshal(v)
}

// Save saves the config to the specified Writer.
func (t *Toml) Save(w io.Writer) (int64, error) {
	contents, err := t.contents()
	if err != nil {
		return 0, err
	}

	n, err := w.Write(contents)
	return int64(n), err
}

// contents returns the config TOML as a byte slice.
// Any required formatting is applied.
func (t Toml) contents() ([]byte, error) {
	content := (*toml.Tree)(&t)

	buffer := bytes.NewBuffer(nil)

	enc := toml.NewEncoder(buffer).Indentation("")
	if err := enc.Encode((*toml.Tree)(content)); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}
	return buffer.Bytes(), nil
}

// Delete deletes the specified key from the TOML config.
func (t *Toml) Delete(key string) error {
	return (*toml.Tree)(t).Delete(key)
}

// Get returns the value for the specified key.
func (t *Toml) Get(key string) interface{} {
	return (*toml.Tree)(t).Get(key)
}

// Set sets the specified key to the specified value in the TOML config.
func (t *Toml) Set(key string, value interface{}) {
	(*toml.Tree)(t).Set(key, value)
}
