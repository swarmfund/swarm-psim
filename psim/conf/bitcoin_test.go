package conf

import (
	"bytes"
	"testing"

	"sync"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func ConfigHelper(t *testing.T, raw string) ViperConfig {
	t.Helper()

	r := bytes.NewReader([]byte(raw))
	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(r)
	if err != nil {
		t.Fatal(err)
	}

	return ViperConfig{
		viper: v,
		Mutex: &sync.Mutex{},
	}
}

func TestViperConfig_Bitcoin(t *testing.T) {
	t.Run("Successfully set config", func(t *testing.T) {
		btcConfigRaw := `
bitcoin:
  node_host: swarm
  node_port: 8332
  node_auth_key: dTAwM2IwNGVmOTUwMWJiZjA2YjpwMDA2NWUyN2MxNWY1NTBiOTJh
  testnet: true
  request_timeout_s: 30
`
		config := ConfigHelper(t, btcConfigRaw)
		assert.NotPanics(t, func() {
			btc := config.Bitcoin()
			assert.NotNil(t, btc)
			assert.Equal(t, true, btc.IsTestnet())
		})

	})

	t.Run("Failed to set bitcoin config", func(t *testing.T) {
		btcConfigRaw := `
bitcoin:
  invalid: 123`
		config := ConfigHelper(t, btcConfigRaw)
		assert.Panics(t, func() {
			btc := config.Bitcoin()
			assert.Nil(t, btc)
		})
	})
}
