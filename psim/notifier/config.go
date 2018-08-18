package notifier

import (
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required"`

	OrderCancelled EventConfig `fig:"order_cancelled,required"`
	KYCCreated     EventConfig `fig:"kyc_created,required"`
	KYCApproved    EventConfig `fig:"kyc_approved,required"`
	USAKyc         EventConfig `fig:"usa_kyc,required"`
	KYCRejected    EventConfig `fig:"kyc_rejected,required"`
	PaymentV2      EventConfig `fig:"payment_v2,required"`
}
