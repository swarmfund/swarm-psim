package providers

import horizon "gitlab.com/swarmfund/horizon-connector"

type Tick interface {
	Ops() []horizon.SetRateOp
}

type Provider interface {
	Errors() chan error
	Ticks() chan Tick
}

type Providers map[string]Provider
