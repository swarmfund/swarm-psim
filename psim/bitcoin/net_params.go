package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

const (
	BTCCurrency  = "BTC"
	DASHCurrency = "DASH"

	MainnetBlockchain = "mainnet"
	TestnetBlockchain = "testnet"
)

func GetNetParams(currency, blockchain string) (*chaincfg.Params, error) {
	switch currency {
	case BTCCurrency:
		return getBTCNetParams(blockchain)
	case DASHCurrency:
		return getDASHNetParams(blockchain)
	default:
		return nil, errors.Errorf("Unsupported currency '%s'.", currency)
	}
}

func getBTCNetParams(blockchain string) (*chaincfg.Params, error) {
	switch blockchain {
	case MainnetBlockchain:
		return &chaincfg.MainNetParams, nil
	case TestnetBlockchain:
		return &chaincfg.TestNet3Params, nil
	default:
		return nil, errors.Errorf("Unsupported blockchain '%s'.", blockchain)
	}
}

func getDASHNetParams(blockchain string) (*chaincfg.Params, error) {
	switch blockchain {
	case MainnetBlockchain:
		return &chaincfg.Params{
			PubKeyHashAddrID: 0x4c,
			ScriptHashAddrID: 0x10,
			PrivateKeyID:     0xcc,
			HDPrivateKeyID:   [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
			HDPublicKeyID:    [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
		}, nil
	case TestnetBlockchain:
		return &chaincfg.Params{
			PubKeyHashAddrID: 0x8c,
			ScriptHashAddrID: 0x13,
			PrivateKeyID:     0xef,
			HDPrivateKeyID:   [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
			HDPublicKeyID:    [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
		}, nil
	default:
		return nil, errors.Errorf("Unsupported blockchain '%s'.", blockchain)
	}
}
