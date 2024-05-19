package config

// RuntimeConfig stores the config options for the container hooks runtime
type RuntimeConfig struct {
	DebugFilePath string `toml:"debug"`
	// LogLevel defines the logging level for the container hooks runtime
	LogLevel string `toml:"log-level"`
	// Runtimes defines the candidates for the low-level runtime
	Runtimes []string `toml:"runtimes"`
}

// GetDefaultRuntimeConfig defines the default values for the config
func GetDefaultRuntimeConfig() (*RuntimeConfig, error) {
	cfg, err := GetDefault()
	if err != nil {
		return nil, err
	}

	return &cfg.ContainerHookRuntimeConfig, nil
}
