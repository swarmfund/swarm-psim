package bearer

import (
	"time"

	"gitlab.com/swarmfund/go/keypair"
)

// Config is structure to parse config for bearer Service into.
type Config struct {
	// Signer is seed of the Master Account Signer,
	// which can create and submit operations.
	Signer *keypair.Full
	// Source is address of the Master Account.
	Source keypair.KP
	// Period is the time duration between the submissions of the operations.
	Period time.Duration
}
