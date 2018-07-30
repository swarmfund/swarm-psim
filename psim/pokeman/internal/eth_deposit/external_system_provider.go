package eth_deposit

import (
	"gitlab.com/swarmfund/psim/psim/internal"
			)

type ExternalSystemProvider interface {
	GetExternalSystemType() (int32, error)
}

type externalSystemProvider struct {
	assetsQ internal.AssetsQ
	asset string
}

func NewExternalSystemProvider(assetsQ internal.AssetsQ, asset string) *externalSystemProvider {
	return &externalSystemProvider{
		assetsQ,
		asset,
	}
}

func (e *externalSystemProvider) GetExternalSystemType() (int32, error) {
	return internal.GetExternalSystemType(e.assetsQ, e.asset)
}