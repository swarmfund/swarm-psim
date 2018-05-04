package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_Ethereum(t *testing.T) {
	ethConfigRaw := `
ethereum:
  proto: http
  host: swarm
  port: 8545`

	config := ConfigHelper(t, ethConfigRaw)
	assert.NotPanics(t, func() {
		eth := config.Ethereum()
		assert.NotNil(t, eth)
	})
}
