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
	// AbnormalPeriod is the time duration between the submissions of
	// the operations, if previous failed.
	AbnormalPeriod time.Duration `fig:"abnormal_period"`
	// SleepPeriod is the time duration between the submissions of
	// the check sales op, if previous no sales found.
	SleepPeriod time.Duration `fig:"sleep_period"`
}
