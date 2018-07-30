package eth_deposit

import (
	"gitlab.com/tokend/horizon-connector"
)

type ExternalBindingDataProvider interface {
	CurrentExternalBindingData() (*string, error)
}

type ebdProvider struct {
	connector *horizon.Connector
	address string
	externalSystem int32
}

func NewExternalBindingDataProvider(connector *horizon.Connector, address string, externalSystem int32) ExternalBindingDataProvider {
	return &ebdProvider{
		connector,
		address,
		externalSystem,
	}
}

func (e *ebdProvider) CurrentExternalBindingData() (*string, error) {
	return e.connector.Accounts().CurrentExternalBindingData(e.address, e.externalSystem)
}
