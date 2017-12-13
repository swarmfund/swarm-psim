package ratesync

import (
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/psim/psim/ratesync/noor"
)

type Config struct {
	ServiceName   string
	LeadershipKey string
	Host          string
	Port          int
	Assets        []Asset
	Pprof         bool
	Signer        keypair.KP
	Master        keypair.KP
	Provider      string
	Noor          NoorConfig
}

type NoorConfig struct {
	Host  string
	Port  int
	Pairs []noor.Pair
}

type Asset struct {
	Code   string
	Ticker string
	Hub    string
	Quote  string
	ToCoin float64 `mapstructure:"to_coin"`
}
