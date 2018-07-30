package eth_deposit

import (
	"context"

	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ExternalSystemBinder interface {
	Bind() error
}

type externalSystemBinder struct {
	builder *xdrbuild.Builder
	connector *horizon.Connector
	source keypair.Address
	signer keypair.Full
	externalSystem int32
}

func NewExternalSystemBinder(builder *xdrbuild.Builder, connector *horizon.Connector, source keypair.Address, signer keypair.Full, externalSystem int32) ExternalSystemBinder {
	return &externalSystemBinder{
		builder,
		connector,
		source,
		signer,
		externalSystem,
	}
}

func (e *externalSystemBinder) Bind() error {
	envelope, err := e.builder.Transaction(e.source).Op(
		&xdrbuild.BindExternalSystemAccountIDOp{e.externalSystem},
	).Sign(e.signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal bind tx")
	}

	result := e.connector.Submitter().Submit(context.Background(), envelope)
	if result.Err != nil {
		return errors.Wrap(result.Err, "failed to submit bind tx", result.GetLoganFields())
	}
	return nil
}