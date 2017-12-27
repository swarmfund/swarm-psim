package horizon

import (
	"fmt"
	"net/url"
	"testing"

	"gitlab.com/swarmfund/horizon-connector/v2"
)

func TestIntegration(t *testing.T) {
	base, _ := url.Parse("http://localhost:8000")
	connector := horizon.NewConnector(base)
	//for {
	//	asset, err := connector.Assets().ByCode("SUN")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println(asset)
	//}

	txs, meta, err := connector.Transactions().Transactions("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(meta)
	fmt.Println(txs)

	events := make(chan horizon.TransactionEvent)
	errs := connector.Listener().Transactions(events)
	for {
		select {
		case err := <-errs:
			fmt.Println(err)
		case event := <-events:
			fmt.Println(event.Meta)
			if event.Transaction != nil {
				fmt.Println(event.Transaction.PagingToken)
			}
		}
	}
	//address := "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
	//account, err := connector.Accounts().ByAddress(address)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(account)
	//
	//address = "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
	//signers, err := connector.Accounts().Signers(address)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(signers)
	//
	//address = "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
	//balances, err := connector.Accounts().Balances(address)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(balances)
	//
	//requests, err := connector.Operations().Requests("")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(requests)
}
