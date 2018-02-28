package pricesetterveri

import (
	"net"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/verifier"
	"gitlab.com/tokend/keypair"
)

func New(
	serviceName string,
	log *logan.Entry,
	config Config,
	pFinder priceFinder,
	signer keypair.Full,
	builder *xdrbuild.Builder,
	listener net.Listener,
	discoveryClient *discovery.Client) app.Service {

	v := newVerifier(
		serviceName,
		log,
		config,
		pFinder,
	)

	return verifier.New(
		serviceName,
		"my_awesome_super_duper_random_id_price_setter",
		log,
		v,
		builder,
		signer,
		listener,
		discoveryClient)
}
