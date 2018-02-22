package airdrop

import "gitlab.com/tokend/keypair"

type Config struct {
	Source keypair.Address `fig:"source"`
	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`
}
