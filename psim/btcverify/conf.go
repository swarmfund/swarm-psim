package btcverify

import "gitlab.com/swarmfund/go/keypair"

type Config struct {
	Host        string
	Port        int
	ServiceName string
	Signer      keypair.KP
	Pprof       bool
}
