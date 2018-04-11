package depositveri

import (
	"net"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/swarmfund/psim/psim/verifier"
	"gitlab.com/tokend/keypair"
)

func New(
	externalSystem string,
	serviceName string,
	log *logan.Entry,
	signer keypair.Full,
	lastBlocksNotWatch uint64,
	// TODO Interface
	horizon *horizon.Connector,
	builder *xdrbuild.Builder,
	listener net.Listener,
	discoveryClient *discovery.Client,
	offchainHelper deposit.OffchainHelper) app.Service {

	v := newVerifier(
		serviceName,
		externalSystem,
		log,
		lastBlocksNotWatch,
		horizon,
		offchainHelper,
	)

	return verifier.New(
		serviceName,
		"my_awesome_super_duper_random_id_deposit",
		log,
		v,
		builder,
		signer,
		listener,
		discoveryClient)
}
