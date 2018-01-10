package horizon

import (
	"testing"
)

func TestNewSubmitError(t *testing.T) {
	//body := []byte(`{
	//	  "type": "https://stellar.org/horizon-errors/transaction_failed",
	//	  "title": "Transaction Failed",
	//	  "status": 400,
	//	  "detail": "The transaction failed when submitted to the stellar network",
	//	  "instance": "ubuntu-8gb-sgp1-01/9o3F5VXjJE-001383",
	//	  "extras": {
	//	    "envelope_xdr": "AAAAAMqNuJ28R4ygz4/7XhWlbWN38eEbvPvcSDaonAaVucJ9AAAAAN7gnfwAAAAAAAAAAAAAAABZafexAAAAAAAAAAEAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAACcQAAAAAJyfGVgnBpcmiXU74TaP2i4/DPWr7WcMkcGSdcfV60pDAAAABFhVU0QAAAA4R0JaQlhJQVEzUEQ2U0JWUEgyU0JJSVBDQU5TMlNLQk5QR1lNNkRLSkczRjVKTFozMldBNDNYSEkAAAAAAAAAAlkeBpAAAABAcelMix939I5EKnn6zCjKqIzofFAotxWd0JND4kJRaVZEy1c+/XKUmjKNPnOGDULvIkxBkdiGQsWIFxv5rrQaDf3onmMAAABAq709ui0FuuJFC6KnuIiI0micPmJ1365paiFjKxxxd4AJvEF4kia50SN6CWFzssHoB2ddhk097nbV5pcBuakrAg==",
	//	    "result_codes": {
	//	      "transaction": "tx_failed",
	//	      "operations": [
	//	        "op_asset_not_found"
	//	      ]
	//	    },
	//	    "result_xdr": "AAAAAAAAAAD/////AAAAAQAAAAAAAAAD////+wAAAAA="
	//	  }
	//	}`)
	//
	//serr, err := NewSubmitError(&http.Response{
	//	Body: ioutil.NopCloser(bytes.NewReader(body)),
	//})
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//txCode := serr.TransactionCode()
	//expectedTxCode := "tx_failed"
	//if txCode != expectedTxCode {
	//	t.Fatalf("expected %s got %s", expectedTxCode, txCode)
	//}
	//
	//opCodes := serr.OperationCodes()
	//expectedOpCodes := []string{"op_asset_not_found"}
	//if !reflect.DeepEqual(opCodes, expectedOpCodes) {
	//	t.Fatalf("expected %s got %s", expectedOpCodes, opCodes)
	//}
	t.Fatal("not implemented")
}
