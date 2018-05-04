package btcwithdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/keypair"
)

func TestNewConfig(t *testing.T) {
	t.Run("Successfully create config", func(t *testing.T) {
		configData := map[string]interface{}{
			"btc_private_key":           "cVaoN6CiScaWB8KJCP8cBP3CjEgFfdK6tf3N5gDow499Xr5AZB8N",
			"hot_wallet_address":        "2N1w4RzejEWkCyumsZY8prvmRxAPFkcwehb",
			"hot_wallet_script_pub_key": "a9145f49aacdc4f9a50e71073e8ed3c449a27759517687",
			"hot_wallet_redeem_script":  "522102cff9f17973e0b1d3468ae29532156f43e42d213fa85e1df40154d7f5748fab6221037afc702c97360f5bd534e6e7eeec0963fd71f9e873e31720ba200c131cfc1f1152ae",
			"min_withdraw_amount":       10000,
			"signer":                    "SAJMOHVPENU2JPK34VR5MXH72LZZX2TABAJLHY5RQ5CXG6XVCMZU3I3N",
		}

		expected := Config{
			PrivateKey:            "cVaoN6CiScaWB8KJCP8cBP3CjEgFfdK6tf3N5gDow499Xr5AZB8N",
			HotWalletAddress:      "2N1w4RzejEWkCyumsZY8prvmRxAPFkcwehb",
			HotWalletScriptPubKey: "a9145f49aacdc4f9a50e71073e8ed3c449a27759517687",
			HotWalletRedeemScript: "522102cff9f17973e0b1d3468ae29532156f43e42d213fa85e1df40154d7f5748fab6221037afc702c97360f5bd534e6e7eeec0963fd71f9e873e31720ba200c131cfc1f1152ae",
			MinWithdrawAmount:     10000,
			SignerKP:              keypair.MustParseSeed("SAJMOHVPENU2JPK34VR5MXH72LZZX2TABAJLHY5RQ5CXG6XVCMZU3I3N"),
		}

		got, err := NewConfig(configData)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, *got)
	})

	t.Run("Failed to set config, because of not enough data", func(t *testing.T) {
		configData := map[string]interface{}{}
		_, err := NewConfig(configData)
		assert.Error(t, err)
	})
}
