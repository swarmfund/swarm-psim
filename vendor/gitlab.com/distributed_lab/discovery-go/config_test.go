package discovery

import "testing"

func TestClientConfig_WithDefaults(t *testing.T) {
	config := ClientConfig{}
	defaults := config.WithDefaults()
	if config.Host != "" {
		t.Error("host should be empty, got %s", config.Host)
	}
	if config.Port != 0 {
		t.Error("port should be 0, got %d", config.Port)
	}
	if defaults.Host != DefaultHost {
		t.Error("host should be %s, got %s", DefaultHost, defaults.Host)
	}
	if defaults.Port != DefaultPort {
		t.Error("port should be %d, got %d", DefaultPort, defaults.Port)
	}

}
