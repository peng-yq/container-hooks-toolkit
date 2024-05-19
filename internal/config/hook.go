package config

// RuntimeHookConfig stores the hooks path 
type RuntimeHookConfig struct {
	// Path specifies the path to the container hooks injected
	Path string `toml:"path"`
}

// GetDefaultRuntimeHookConfig defines the default values for the config
func GetDefaultRuntimeHookConfig() (*RuntimeHookConfig, error) {
	cfg, err := GetDefault()
	if err != nil {
		return nil, err
	}

	return &cfg.ContainerRuntimeHookConfig, nil
}
