package notifier

import (
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`

	SaleCancelled EmailsConfig `fig:"sale_cancelled"`
	KYCCreated    EmailsConfig `fig:"kyc_created"`
	KYCApproved   EmailsConfig `fig:"kyc_approved"`
	KYCRejected   EmailsConfig `fig:"kyc_rejected"`
}
