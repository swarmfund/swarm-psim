package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperConfig_Notificator(t *testing.T) {
	notificatorConfigRaw := `
notificator:
  url: http://swarm:9009
  secret: GB2JMWBAUWR4ZMI4XCA3NYS6PE7AQKINQGCB6M4SXQBUB4LXKPDFCLHX
  public: SA3RNFJWYW4TF5ZMCU54QWYNRP5GOUUDZPT26UYX357L3GNKJ3IJRCFR
`
	config := ConfigHelper(t, notificatorConfigRaw)

	assert.NotPanics(t, func() {
		notificator := config.Notificator()
		assert.NotNil(t, notificator)
	})
}

func TestViperConfig_NotificatorFailed(t *testing.T) {
	notificatorConfigRaw := `
notificator:
  invalid: 123`

	config := ConfigHelper(t, notificatorConfigRaw)
	assert.Panics(t, func() {
		notificator := config.Notificator()
		assert.Nil(t, notificator)
	})
}
