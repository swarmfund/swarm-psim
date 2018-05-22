package listener

import (
	"gitlab.com/tokend/keypair"
)

// Config is structure to parse config for listener Service into.
type Config struct {
	Signer keypair.Full
}
