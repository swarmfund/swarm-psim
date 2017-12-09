package discovery

type ClientConfig struct {
	Env  string
	Host string
	Port int
}

// WithDefaults populates fields with default values if need, returns mutated
// config, does not mutate original, safe to call on nil pointers
func (old *ClientConfig) WithDefaults() *ClientConfig {
	var new ClientConfig
	if old != nil {
		new = *old
	}

	if new.Host == "" {
		new.Host = DefaultHost
	}

	if new.Port == 0 {
		new.Port = DefaultPort
	}

	if new.Env == "" {
		new.Env = DefaultEnv
	}

	return &new
}
