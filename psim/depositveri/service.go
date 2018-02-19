package depositveri

import (
	"net"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/deposit"
	"gitlab.com/tokend/keypair"
	"gitlab.com/swarmfund/psim/psim/verify"
	"gitlab.com/swarmfund/psim/psim/app"
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

	verifier := newVerifier(
		serviceName,
		externalSystem,
		log,
		lastBlocksNotWatch,
		horizon,
		offchainHelper,
	)

	return verify.New(serviceName,
		"my_awesome_super_duper_random_id_deposit",
			log, verifier, builder, signer, listener, discoveryClient)
}
