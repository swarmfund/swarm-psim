package eth

import (
	"testing"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/keypair"
)

func TestNewWithdrawConfig(t *testing.T) {
	t.Run("Success to create config eth withdrawals", func(t *testing.T) {

		configData := map[string]interface{}{
			"asset":                 "ETH",
			"threshold":             123,
			"key":                   "53414c4644494f3447454f54494b4e4743344d4b553345565954354d58444535564343525951373653455a5a4d4e464a584935454b33484f",
			"gas_price":             1000000000,
			"signer":                "SCBDFODTCFIXMC4J634W7UT4NXFN5KNUJGWY3UJ5GISZE4XUGXG4JG6X",
			"token":                 "0xAeEc6f5aCA72F3A005af1B3420ab8c8c7009BaC8",
			"verifier_service_name": "eth_withdraw_verify",
		}

		token := common.Address([20]byte{174, 236, 111, 90, 202, 114, 243, 160, 5, 175, 27, 52, 32, 171, 140, 140, 112, 9, 186, 200})

		expected := WithdrawConfig{
			Signer:              keypair.MustParseSeed("SCBDFODTCFIXMC4J634W7UT4NXFN5KNUJGWY3UJ5GISZE4XUGXG4JG6X"),
			Asset:               "ETH",
			Threshold:           123,
			Token:               &token,
			Key:                 "53414c4644494f3447454f54494b4e4743344d4b553345565954354d58444535564343525951373653455a5a4d4e464a584935454b33484f",
			GasPrice:            big.NewInt(1000000000),
			VerifierServiceName: "eth_withdraw_verify",
		}

		got, err := NewWithdrawConfig(configData)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, *got)
	})

	t.Run("Failed to create config eth withdrawals, because of not enough data", func(t *testing.T) {
		configData := map[string]interface{}{}
		_, err := NewWithdrawConfig(configData)
		assert.Error(t, err)
	})

}
