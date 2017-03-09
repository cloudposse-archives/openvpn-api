package config

// Config - structure to store global configuration
type Config struct {
	Listen string
}

// Validate - process validation of config values
func (c *Config) Validate() (err error) {
	err = nil
	return
}
