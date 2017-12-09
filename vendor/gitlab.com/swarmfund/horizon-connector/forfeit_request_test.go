package horizon

import (
	"encoding/json"
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal(fmt.Sprintf("%v of %T != %v of %T", a, a, b, b))
}

func TestForfeitRequestStruct(t *testing.T) {
	data := []byte(`{
                "accepted": null,
                "created_at": "2017-06-29T19:56:13Z",
                "exchange": "GDFI3OE5XRDYZIGPR75V4FNFNVRXP4PBDO6PXXCIG2UJYBUVXHBH2XVR",
                "paging_token": "1",
                "payment_details": {
                    "amount": "100.0000",
                    "asset": "",
                    "fee_from_source": false,
                    "fixed_fee": "0.0000",
                    "from": "GBUI5HXG3IFXZYT6UJFOP6IEYNRTENOLAQJ7KAXBBFDUUE5BCTY6LNMO",
                    "from_balance": "BDPOP4TDCQUEBLDP5DDR26EOTM333SEDO6DO5DMSTR37LBNSRNBVAKGX",
                    "payment_fee": "0.0000",
                    "to": "",
                    "to_balance": "",
                    "user_details": "foo@bar.com"
                },
                "payment_id": "6",
                "payment_state": 1,
                "request_type": 2,
                "updated_at": "2017-06-29T19:56:13Z"
            }`)

	var f ForfeitRequest
	err := json.Unmarshal(data, &f)
	if err != nil {
		t.Fatal(err)
	}

	var null *bool
	assertEqual(t, f.Accepted, null)
	assertEqual(t, f.PaymentID, "6")
	assertEqual(t, f.PaymentDetails.Amount, "100.0000")
	assertEqual(t, f.ID, "1")
}
