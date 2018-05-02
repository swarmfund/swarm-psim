package pricesetterveri

import (
	"net"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/utils"
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

	v := NewVerifier(serviceName, log, config, pFinder)

	return verifier.New(
		serviceName,
		utils.GenerateToken(),
		log,
		v,
		builder,
		signer,
		listener,
		discoveryClient)
}
