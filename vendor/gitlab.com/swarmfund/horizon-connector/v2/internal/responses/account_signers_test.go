package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountSignersUnmarshal(t *testing.T) {
	data := []byte(`{
    "signers": [
        {
            "public_key": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
            "signer_identity": 1,
            "signer_name": "foobar",
            "signer_type": {
                "flags": [
                    {
                        "name": "reader",
                        "value": 1
                    },
                    {
                        "name": "not_verified_acc_manager",
                        "value": 2
                    },
                    {
                        "name": "general_acc_manager",
                        "value": 4
                    },
                    {
                        "name": "direct_debit_operator",
                        "value": 8
                    },
                    {
                        "name": "asset_manager",
                        "value": 16
                    },
                    {
                        "name": "asset_rate_manager",
                        "value": 32
                    },
                    {
                        "name": "balance_manager",
                        "value": 64
                    },
                    {
                        "name": "issuance_manager",
                        "value": 128
                    },
                    {
                        "name": "invoice_manager",
                        "value": 256
                    },
                    {
                        "name": "payment_operator",
                        "value": 512
                    },
                    {
                        "name": "limits_manager",
                        "value": 1024
                    },
                    {
                        "name": "account_manager",
                        "value": 2048
                    },
                    {
                        "name": "commission_balance_manager",
                        "value": 4096
                    },
                    {
                        "name": "operational_balance_manager",
                        "value": 8192
                    }
                ],
                "int": 16383
            },
            "weight": 1
        }
    ]}`)

	var got AccountSigners
	if err := json.Unmarshal(data, &got); err != nil {
		assert.NoError(t, err)
	}

	assert.Len(t, got.Signers, 1)
	signer := got.Signers[0]
	assert.EqualValues(t, 1, signer.Weight)
	assert.EqualValues(t, 16383, signer.Type)
	assert.EqualValues(t, "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636", signer.PublicKey)
	assert.EqualValues(t, 1, signer.Identity)
	assert.EqualValues(t, "foobar", signer.Name)
}
