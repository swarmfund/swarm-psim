package eth

import (
	"encoding/hex"

	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip32"
	"sync"
	"context"
)

var (
	ErrNoKey = errors.New("wallet doesn't have requested key")
)

type Wallet struct {
	hd     bool
	master *bip32.Key
	keys   map[common.Address]ecdsa.PrivateKey

	extendWait sync.WaitGroup
}

func NewHDWallet(hexseed string) (*Wallet, error) {
	seed, err := hex.DecodeString(hexseed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode seed")
	}

	master, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init master")
	}
	master, err = master.NewChildKey(bip32.FirstHardenedChild)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive m/0`")
	}
	master, err = master.NewChildKey(bip32.FirstHardenedChild)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive m/0`/0`")
	}

	wallet := &Wallet{
		hd:     true,
		master: master,
		keys:   make(map[common.Address]ecdsa.PrivateKey),
	}

	wallet.extendWait.Add(1)

	// TODO check horizon for account sequence and extended keys as needed
	if err := wallet.extend(2 << 10); err != nil {
		return nil, errors.Wrap(err, "failed to extend master")
	}

	return wallet, nil
}

func NewWallet() *Wallet {
	return &Wallet{
		keys: make(map[common.Address]ecdsa.PrivateKey),
	}
}

func (wallet *Wallet) ImportHEX(data string) (common.Address, error) {
	raw, err := hex.DecodeString(data)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to decode string")
	}
	return wallet.Import(raw)
}

func (wallet *Wallet) Import(raw []byte) (common.Address, error) {
	pk, err := crypto.ToECDSA(raw)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to convert pk")
	}
	address := crypto.PubkeyToAddress(pk.PublicKey)
	wallet.keys[address] = *pk
	return address, nil
}

func (wallet *Wallet) extend(i uint) error {
	for uint(len(wallet.keys)) < i {
		child, err := wallet.master.NewChildKey(uint32(len(wallet.keys)))
		if err != nil {
			return errors.Wrap(err, "failed to extend child")
		}

		if _, err := wallet.Import(child.Key); err != nil {
			return errors.Wrap(err, "failed to import key")
		}
	}

	wallet.extendWait.Done()

	return nil
}

// TODO Use ctx
func (wallet *Wallet) Addresses(ctx context.Context) (result []common.Address) {
	wallet.extendWait.Wait()

	for addr := range wallet.keys {
		result = append(result, addr)
	}
	return result
}

// TODO pass ctx
func (wallet *Wallet) SignTX(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	wallet.extendWait.Wait()

	key, ok := wallet.keys[address]
	if !ok {
		return nil, ErrNoKey
	}
	return wallet.SignTXWithPrivate(&key, tx)
}

func (wallet *Wallet) SignTXWithPrivate(key *ecdsa.PrivateKey, tx *types.Transaction) (*types.Transaction, error) {
	return types.SignTx(tx, types.HomesteadSigner{}, key)
}

// TODO pass ctx
func (wallet *Wallet) HasAddress(address common.Address) bool {
	wallet.extendWait.Wait()

	_, ok := wallet.keys[address]
	return ok
}
