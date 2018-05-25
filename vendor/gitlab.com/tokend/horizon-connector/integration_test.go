package horizon_test

import (
	"testing"

	"net/url"

	"context"
	"fmt"

	"encoding/json"
	"math"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

func connector() *horizon.Connector {
	base, _ := url.Parse("https://api.testnet.tokend.org")
	return horizon.NewConnector(base)
}

func builder(horizon *horizon.Connector) *xdrbuild.Builder {
	info, err := connector().System().Info()
	if err != nil {
		panic(errors.Wrap(err, "failed to get info"))
	}
	return xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
}

func TestCreateAccount(t *testing.T) {
	//t.Skip("integration")
	// testnet master
	kp := keypair.MustParseSeed("SBRWM6VULCYS5WHEP6DMWOHSABXILU4CYC44BON3VXVOIZUUXJOCOI5I")
	//account, _ := keypair.Random()
	recovery, _ := keypair.Random()
	envelope, err := builder(connector()).Transaction(kp).Op(&xdrbuild.CreateAccountOp{
		Address:     "GDGQI3SSB7N7YDBGWCZB3DT7SA23KJWDTYQB5HCYR5VP3EBD6CXQXXG4",
		AccountType: xdr.AccountTypeSyndicate,
		Recovery:    recovery.Address(),
	}).Sign(kp).Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal tx"))
	}
	result := connector().Submitter().Submit(context.Background(), envelope)
	fmt.Println(result.Err, result.TXCode, result.OpCodes)
	//fmt.Println(account.Seed(), account.Address())
}

func TestBindExternal(t *testing.T) {
	//t.Skip("integration")
	kp := keypair.MustParseSeed("SAXMBVFBTYCVSCDH3OK3BOSY6N2UVRNKLEKXIPI3S6HIM5WBOJ7DXLBV")
	envelope, err := builder(connector()).Transaction(kp).Op(
		&xdrbuild.BindExternalSystemAccountIDOp{9},
	).Sign(kp).Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal tx"))
	}
	result := connector().Submitter().Submit(context.Background(), envelope)
	fmt.Println(result.Err, result.OpCodes, result.TXCode)
}

// TODO implement tests
type CreateAsset struct {
	Code                   string
	AssetSigner            string
	MaxIssuanceAmount      uint64
	InitialPreissuedAmount uint64
	Policies               uint32
	Details                json.RawMessage
}

func (op *CreateAsset) Validate() error {
	// TODO implement
	return nil
}

func (op *CreateAsset) XDR() (*xdr.Operation, error) {
	var assetSigner xdr.AccountId
	if err := assetSigner.SetAddress(op.AssetSigner); err != nil {
		return nil, errors.Wrap(err, "failed to set asset signer")
	}
	xdrop := &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageAsset,
			ManageAssetOp: &xdr.ManageAssetOp{
				Request: xdr.ManageAssetOpRequest{
					Action: xdr.ManageAssetActionCreateAssetCreationRequest,
					CreateAsset: &xdr.AssetCreationRequest{
						Code:                   xdr.AssetCode(op.Code),
						PreissuedAssetSigner:   assetSigner,
						MaxIssuanceAmount:      xdr.Uint64(op.MaxIssuanceAmount),
						InitialPreissuedAmount: xdr.Uint64(op.InitialPreissuedAmount),
						Policies:               xdr.Uint32(op.Policies),
						Details:                xdr.Longstring(string(op.Details)),
					},
				},
			},
		},
	}
	return xdrop, nil
}

func TestCreateAsset(t *testing.T) {
	//t.Skip("integration")
	kp := keypair.MustParseSeed("SBQ3YVRINQOJDT6FQD3EFMZ5THZFWVSQ37RZYTFZBF3TEE5GPBS6NNXD")
	envelope, err := builder(connector()).Transaction(kp).Op(
		&CreateAsset{
			"SWM",
			kp.Address(),
			math.MaxUint64,
			math.MaxUint64,
			123,
			[]byte(`{}`),
		},
	).Sign(kp).Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal tx"))
	}
	result := connector().Submitter().Submit(context.Background(), envelope)
	fmt.Println(result.Err, result.OpCodes, result.TXCode)

}

func TestGetAccount(t *testing.T) {
	kp := keypair.MustParseSeed("SDXTZ4N2OUKLVXKHDYPMHIBRHR2IFZBUIQOGSL32WX2CSP7QGHYTM4WW")
	fmt.Println(kp.Address())
}

func TestConnectorV2(t *testing.T) {
	//t.Skip("integration")

	// create account
	{

	}

	//kp := keypair.MustParseSeed("SDVTI6SGDAFM6VWAWEP7EKMH26HFMJ4DSYRNANT4RSNTEEB7RBRX4Q47")

	//{
	//	q := connector.Wallets()
	//	{
	//		verified := false
	//		page := int32(2)
	//
	//		ops := types.GetOpts{
	//			Verified: &verified,
	//			Page:     &page,
	//		}
	//
	//		wallets, resultPage, err := q.Filter(&ops)
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//
	//		fmt.Println(wallets)
	//		fmt.Println(resultPage)
	//	}
	//
	//	{
	//		err := q.Delete("22c25cd3151661a01d7f0c502169f4acf4c0d366a606ec4194334f3846b0b195")
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//	}
	//}
	//
	//{
	//	q := connector.Templates()
	//	{
	//		body := strings.NewReader(`<!DOCTYPE html>
	//		<html>
	//		<body>
	//
	//		<h1>My First Heading</h1>
	//
	//		<p>My first paragraph.</p>
	//
	//		</body>
	//		</html>
	//		`)
	//
	//		_, err := q.Put("test", body)
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//	}
	//	{
	//		body, err := q.Get("test")
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//
	//		t.Log(string(body))
	//	}
	//}
	//
	//{
	//	r, err := connector.Client().Post("/participants", strings.NewReader(`{
	//		"for_account": "GDHK26UFBGC63UBQCVQLHJD6RAQXLAS7RKJAR5FZQAWMCUBFHRNKFSKC",
	//		"participants": {"1": [
	//			{
	//				"account_id": "GDS67HI27XJIJEL7IGHVJVNHPXZLMW6F3O45OXIMKAUNGIR2ROBUKTT4"
	//			},
	//			{
	//				"account_id": "GDHK26UFBGC63UBQCVQLHJD6RAQXLAS7RKJAR5FZQAWMCUBFHRNKFSKC"
	//			}
	//		]}
	//	}`))
	//	if err != nil {
	//		herr, ok := err.(Error)
	//		if ok {
	//			fmt.Println(string(herr.Body()))
	//		}
	//		t.Fatal(err)
	//	}
	//	fmt.Println(string(r))
	//}
	//
	//{
	//	assets, err := connector.Assets().Index()
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println(assets)
	//}
	//
	//{
	//	kp, _ := keypair.Random()
	//	kp2, _ := keypair.Random()
	//	envelope, err := xdrbuild.
	//		NewBuilder("Test SDF Network ; September 2015", 3600).
	//		Transaction(keypair.MustParseAddress("GDHK26UFBGC63UBQCVQLHJD6RAQXLAS7RKJAR5FZQAWMCUBFHRNKFSKC")).
	//		Op(xdrbuild.CreateAccountOp{
	//			Address:     kp.Address(),
	//			AccountType: 2,
	//			Recovery:    kp2.Address(),
	//		}).Sign(keypair.MustParseSeed("SB3YDBQV7VPJEWBT5FLSKO5N2WMAFJR46JXPV7HKTANXW4IKTMKZ2VNE")).Marshal()
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	submitter := connector.Submitter()
	//	result := submitter.Submit(context.TODO(), envelope)
	//	fmt.Printf("%#v\n", result)
	//}
	//
	//for {
	//	asset, err := connector.Assets().ByCode("SUN")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println(asset)
	//}
	//
	//txs, meta, err := connector.Transactions().Transactions("")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(meta)
	//fmt.Println(txs)
	//
	//events := make(chan TransactionEvent)
	//errs := connector.Listener().Transactions(events)
	//for {
	//	select {
	//	case err := <-errs:
	//		fmt.Println(err)
	//	case event := <-events:
	//		fmt.Println(event.Meta)
	//		if event.Transaction != nil {
	//			fmt.Println(event.Transaction.PagingToken)
	//		}
	//	}
	//}
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
	////requests, err := connector.Operations().Requests("")
	////if err != nil {
	////	t.Fatal(err)
	////}
	//fmt.Println(requests)
}
