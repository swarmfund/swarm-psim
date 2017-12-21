package resources

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestUnmarshal(t *testing.T) {
	cases := []struct {
		name     string
		data     string
		expected Request
	}{
		{
			"withdraw",
			`{
				"details": {
					"request_type": "withdraw",
					"request_type_i": 4,
					"withdraw": {
						"amount": "10000.0000",
						"balance_id": "BCLS6FR7XLDCCTVRSTDUIJ6LQEHHZKOQNSTUZBDIBBQLIABCTVRLG6QE",
						"dest_asset_amount": "10000.0000",
						"dest_asset_code": "BTC466",
						"external_details": "Random external details",
						"fixed_fee": "0.0000",
						"percent_fee": "0.0000"
					}
				},
				"hash": "14291500540b187bd9266bd937a22e34da8705f768bbb07f5428c968f315e355",
				"id": "14",
				"paging_token": "14",
				"reference": null,
				"reject_reason": "",
				"request_state": "approved",
				"request_state_i": 3,
				"requestor": "GCHF3FNFGO25YCZWXZVHNA3OKVC2CQZ7RPVT2LOFWZFZJAMOZSZI3EOA",
				"reviewer": "GDOWZD2OFT2TOB3VPAVAFMG5RYACWDC6HAMOSZZ7IESAOTUOMMN7IJJS"
			}`,
			Request{
				Hash:        "14291500540b187bd9266bd937a22e34da8705f768bbb07f5428c968f315e355",
				State:       3,
				ID:          14,
				PagingToken: "14",
				Details: RequestDetails{
					RequestType: 4,
					Withdraw: &RequestWithdrawDetails{
						Amount:           100000000,
						BalanceID:        "BCLS6FR7XLDCCTVRSTDUIJ6LQEHHZKOQNSTUZBDIBBQLIABCTVRLG6QE",
						DestinationAsset: "BTC466",
						ExternalDetails:  "Random external details",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got Request
			if err := json.Unmarshal([]byte(tc.data), &got); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.expected, got)
		})
	}
}
