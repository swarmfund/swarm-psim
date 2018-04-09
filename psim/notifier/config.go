package notifier

import (
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	SaleCancelled EventConfig `fig:"sale_cancelled"`
	KYCCreated    EventConfig `fig:"kyc_created"`
	KYCApproved   EventConfig `fig:"kyc_approved"`
	KYCRejected   EventConfig `fig:"kyc_rejected"`
}
