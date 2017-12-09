package stripeverify

import "gitlab.com/tokend/go/keypair"

type Config struct {
	Host        string
	Port        int
	ServiceName string
	Signer      keypair.KP
	Pprof       bool
}
