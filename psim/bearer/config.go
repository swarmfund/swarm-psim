package bearer

import (
	"time"

	"gitlab.com/tokend/keypair"
)

// Config is structure to parse config for bearer Service into.
type Config struct {
	//where to host service
	Host string `fig:"host,required"`
	//which port to open
	Port int `fig:"port,required"`
	// Signer is seed of the Master Account Signer,
	// which can create and submit operations.
	Signer keypair.Full `fig:"signer,required"`
	// Source is address of the Master Account.
	Source     keypair.Address `fig:"source,required"`
	NormalTime time.Duration   `fig:"normal_time,required"`
	// AbnormalPeriod is the time duration between the submissions of
	// the operations, if previous failed.
	AbnormalPeriod time.Duration `fig:"abnormal_period"`
	//MaxAbnormalPeriod is the maximum time duration between the submission
	//of operation if previous failed
	MaxAbnormalPeriod time.Duration `fig:"max_abnormal_period"`
	// SleepPeriod is the time duration between the submissions of
	// the check sales op, if previous no sales found.
}
