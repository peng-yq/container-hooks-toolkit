package config

// CTKConfig stores the config options for container-hooks-ctk
type CTKConfig struct {
	Path string `toml:"path"`
}

// GetDefaultCTKConfig defines the default values for the config
func GetDefaultCTKConfig() (*CTKConfig, error) {
	cfg, err := GetDefault()
	if err != nil {
		return nil, err
	}

	return &cfg.ContainerHookCtkConfig, nil
}
