package ethfunnel

import (
	"testing"

	"encoding/hex"

	"io/ioutil"

	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/piotrnar/gocoin/lib/btc"
)

func TestExtend(t *testing.T) {
	//public := "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB"
	//private := "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"

	seed, err := hex.DecodeString("fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542")
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	ks := keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
	ks = ks
	wallet := btc.MasterKey(seed, false)
	curve := crypto.S256()
	curve = curve

	//0x8d718DD1988Ad38B5576aaEF5785A7abD8902F0C
	for i := uint32(0); i < 10; i++ {
		child := wallet.Child(i)

		//{
		//	pk, err := crypto.ToECDSA(child.Key[1:])
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	bytes := elliptic.Marshal(curve, pk.X, pk.Y)
		//}

		//
		//{
		//	// addr from public
		//	address := crypto.Keccak256Hash(child.Pub().Key[1:])
		//	fmt.Printf("%d %x \n", i, address[12:])
		//}
		//
		{
			//	//pk, err := crypto.ToECDSA(crypto.Keccak256(child.Key[1:]))
			//	//pk, err := crypto.ToECDSA(crypto.Keccak256(child.Key))
			pk, err := crypto.ToECDSA(child.Key[1:])
			if err != nil {
				t.Fatal(err)
			}
			pk = pk

			//	fmt.Println("P", child.Pub().Key)
			//	fmt.Println("M", elliptic.Marshal(curve, pk.X, pk.Y))
			//
			account, err := ks.ImportECDSA(pk, "")
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(account.Address.Hex())
		}
	}
}
