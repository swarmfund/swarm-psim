package balancereporter

import "gitlab.com/tokend/keypair"

type ServiceConfig struct {
	Signer    keypair.Full `fig:"signer,required"`
	AssetCode string       `fig:"asset_code,required"`
}
