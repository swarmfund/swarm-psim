package btcwithdveri

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/conf"
	"net"
)

type BTCClient interface {
	SignAllTXInputs(txHex, scriptPubKey string, redeemScript *string, privateKey string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
	IsTestnet() bool
}

type Service struct {
	log    *logan.Entry
	config Config

	horizon   *horizon.Connector
	btcClient BTCClient
	listener  net.Listener
}

func New(log *logan.Entry, config Config, horizon *horizon.Connector, btc BTCClient, listener net.Listener) *Service {

	return &Service{
		log:    log.WithField("service", conf.ServiceBTCWithdrawVerify),
		config: config,

		horizon:   horizon,
		btcClient: btc,
		listener:  listener,
	}
}

func (s *Service) Run(ctx context.Context) chan error {
	s.serveAPI(ctx)

	errs := make(chan error)
	close(errs)
	return errs
}
