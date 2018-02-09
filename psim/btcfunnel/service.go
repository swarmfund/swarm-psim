package btcfunnel

import (
	"context"

	"encoding/hex"

	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

const (
	outChanSize = 1000
)

// BTCClient is the interface to be implemented by a
// Bitcoin client to parametrize the Service.
type BTCClient interface {
	GetNetParams() *chaincfg.Params
	EstimateFee() (float64, error)

	GetBlockCount() (uint64, error)
	GetBlock(blockNumber uint64) (*btcutil.Block, error)

	GetTxUTXO(txHash string, outNumber uint32) (*bitcoin.UTXO, error)
	GetAddrUTXOs(address string) ([]bitcoin.WalletUTXO, error)

	CreateRawTX(inputUTXOs []bitcoin.Out, addrToAmount map[string]float64) (resultTXHex string, err error)
	SignRawTX(initialTXHex string, inputUTXOs []bitcoin.InputUTXO, privateKeys []string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
}

type UTXO struct {
	bitcoin.UTXO
	bitcoin.Out
}

// Service implements app.Service to be registered in the app.
type Service struct {
	config Config
	log    *logan.Entry

	lastProcessedBlock uint64
	addrToPriv         map[string]string

	btcClient BTCClient
}

// New is constructor for btcfunnel Service.
func New(config Config, log *logan.Entry, btcClient BTCClient) *Service {
	return &Service{
		config: config,
		log:    log,

		lastProcessedBlock: config.LastProcessedBlock,
		addrToPriv:         make(map[string]string),

		btcClient: btcClient,
	}
}

// Run is implementation of app.Service, Run is called by the app.
// Run will return only when work is finished.
func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	err := s.deriveKeys()
	if err != nil {
		// Don't try again, because Keys derivation process does not depend on anything, so if failed once - will fail always.
		s.log.WithError(err).Error("Failed to derive keys from the extended private key from config, stopping.")
		return
	}

	app.RunUntilSuccess(ctx, s.log, "existing_blocks_fetcher", s.fetchExistingBlocks, 5*time.Second)
	if app.IsCanceled(ctx) {
		return
	}

	s.log.Info("Started listening to newly appeared Blocks.")
	app.RunOverIncrementalTimer(ctx, s.log, "new_blocks_fetcher", s.fetchNewBlock, 10*time.Second, 5*time.Second)
}

func (s *Service) deriveKeys() error {
	s.log.WithField("keys_to_derive", s.config.KeysToDerive).Info("Started keys deriving.")

	extKey, err := hdkeychain.NewKeyFromString(s.config.ExtendedPrivateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to create new ExtendedPrivateKey from the string representation")
	}

	extKey.SetNet(s.btcClient.GetNetParams())

	for i := uint64(0); i < s.config.KeysToDerive; i++ {
		fields := logan.F{
			"child_i": i,
		}

		child, err := extKey.Child(uint32(i))
		if err != nil {
			return errors.Wrap(err, "Failed to derive Child of the ExtendedPrivateKey", fields)
		}

		priv, err := child.ECPrivKey()
		if err != nil {
			return errors.Wrap(err, "Failed to get ECPrivKey from the Child", fields)
		}

		wif, err := btcutil.NewWIF(priv, s.btcClient.GetNetParams(), true)
		if err != nil {
			return errors.Wrap(err, "Failed to get WIF key from the PrivKey", fields)
		}

		pubKeyHash := btcutil.Hash160(priv.PubKey().SerializeCompressed())
		addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, s.btcClient.GetNetParams())
		if err != nil {
			return errors.Wrap(err, "Failed to create P2PKH Address of the PrivKey", fields.Merge(logan.F{
				"pub_key_hash": hex.EncodeToString(pubKeyHash),
			}))
		}

		s.addrToPriv[addr.String()] = wif.String()
	}

	s.log.WithField("keys_to_derive", s.config.KeysToDerive).Info("Finished keys deriving.")
	return nil
}
