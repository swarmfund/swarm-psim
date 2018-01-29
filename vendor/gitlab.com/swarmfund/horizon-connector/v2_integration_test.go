package horizon

import (
	"fmt"
	"net/url"
	"testing"

	"strings"

	"context"

	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/tokend/keypair"
)

func TestConnectorV2(t *testing.T) {
	t.Skip("integration")

	base, _ := url.Parse("http://dev.swarm:8000")
	master := keypair.MustParseSeed("SB3YDBQV7VPJEWBT5FLSKO5N2WMAFJR46JXPV7HKTANXW4IKTMKZ2VNE")
	connector := horizon.NewConnector(base).WithSigner(master)

	{
		r, err := connector.Client().Post("/participants", strings.NewReader(`{
			"for_account": "GDHK26UFBGC63UBQCVQLHJD6RAQXLAS7RKJAR5FZQAWMCUBFHRNKFSKC",
			"participants": {"1": [
				{
					"account_id": "GDS67HI27XJIJEL7IGHVJVNHPXZLMW6F3O45OXIMKAUNGIR2ROBUKTT4"
				}, 
				{
					"account_id": "GDHK26UFBGC63UBQCVQLHJD6RAQXLAS7RKJAR5FZQAWMCUBFHRNKFSKC"
				}
			]}
		}`))
		if err != nil {
			herr, ok := err.(horizon.Error)
			if ok {
				fmt.Println(string(herr.Body()))
			}
			t.Fatal(err)
		}
		fmt.Println(string(r))
	}

	{
		assets, err := connector.Assets().Index()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(assets)
	}

	{
		kp, _ := keypair.Random()
		kp2, _ := keypair.Random()
		envelope, err := xdrbuild.
			NewBuilder("Test SDF Network ; September 2015", 3600).
			Transaction(keypair.MustParseAddress("GDHK26UFBGC63UBQCVQLHJD6RAQXLAS7RKJAR5FZQAWMCUBFHRNKFSKC")).
			Op(xdrbuild.CreateAccountOp{
				Address:     kp.Address(),
				AccountType: 2,
				Recovery:    kp2.Address(),
			}).Sign(keypair.MustParseSeed("SB3YDBQV7VPJEWBT5FLSKO5N2WMAFJR46JXPV7HKTANXW4IKTMKZ2VNE")).Marshal()
		if err != nil {
			t.Fatal(err)
		}
		submitter := connector.Submitter()
		result := submitter.Submit(context.TODO(), envelope)
		fmt.Printf("%#v\n", result)
	}

	for {
		asset, err := connector.Assets().ByCode("SUN")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(asset)
	}

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
	address := "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
	account, err := connector.Accounts().ByAddress(address)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(account)

	address = "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
	signers, err := connector.Accounts().Signers(address)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(signers)

	address = "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
	balances, err := connector.Accounts().Balances(address)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balances)

	requests, err := connector.Operations().Requests("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(requests)
}
