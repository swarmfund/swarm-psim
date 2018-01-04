package bearer

import (
	"time"

	"gitlab.com/swarmfund/go/keypair"
)

type Config struct {
	Signer *keypair.Full
	Period time.Duration
}
