package ethfunnel

import (
	"testing"

	"encoding/hex"

	"io/ioutil"

	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/piotrnar/gocoin/lib/btc"
	"github.com/stretchr/testify/assert"
)

func TestExtend(t *testing.T) {
	seed, err := hex.DecodeString("fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542")
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)

	ks := keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
	wallet := btc.MasterKey(seed, false)

	child := wallet.Child(1)

	pk, err := crypto.ToECDSA(child.Key[1:])
	if err != nil {
		t.Fatal(err)
	}

	account, err := ks.ImportECDSA(pk, "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, account.Address.Hex(), "0xcBA17b04AE211B4f3f3DFf1266e0A8910AC9e3e3")
}
