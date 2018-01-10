package horizon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/strkey"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/internal/resources"
)

type Connector struct {
	baseURL string
	client  *http.Client
	info    *Info
}

func NewConnector(endpoint string) (*Connector, error) {
	info, err := NewInfo(endpoint)
	if err != nil {
		return nil, err
	}

	return &Connector{
		baseURL: strings.TrimRight(endpoint, "/"),
		client:  &http.Client{},
		info:    info,
	}, nil
}

func (c *Connector) Info() (*Info, error) {
	return NewInfo(c.baseURL)
}

func (c *Connector) TimeBounds() xdr.TimeBounds {
	var fixToRemove int64 = 60
	return xdr.TimeBounds{
		MaxTime: xdr.Uint64(time.Now().Unix() + c.info.TxExpirationPeriod - fixToRemove),
	}
}

func (c *Connector) NewBalanceID() (*xdr.BalanceId, error) {
	kp, err := keypair.Random()
	if err != nil {
		return nil, err
	}

	raw, err := strkey.Decode(strkey.VersionByteAccountID, kp.Address())
	if err != nil {
		return nil, err
	}

	var ui xdr.Uint256
	copy(ui[:], raw)
	bid, err := xdr.NewBalanceId(xdr.CryptoKeyTypeKeyTypeEd25519, ui)
	if err != nil {
		return nil, err
	}

	return &bid, nil
}

func (c *Connector) doSigned(kp keypair.KP, method, path string) (*http.Response, error) {
	req, err := NewSignedRequest(c.baseURL, method, path, kp)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *Connector) do(method, path string) (*http.Response, error) {
	req, err := NewRequest(c.baseURL, method, path)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func prepareQueryString(params interface{}) (url.Values, error) {
	// going struct -> json -> map -> form -> string
	paramsBytes, err := json.Marshal(&params)
	if err != nil {
		return nil, err
	}

	paramsMap := map[string]interface{}{}
	err = json.Unmarshal(paramsBytes, &paramsMap)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	for key, value := range paramsMap {
		// TODO type cast non-strings
		switch v := value.(type) {
		case string:
			form.Add(key, v)
		default:
			return nil, errors.New("invalid param type")
		}
	}

	return form, nil
}

func (c *Connector) SignedRequest(method, endpoint string, kp keypair.KP) (*http.Request, error) {
	return NewSignedRequest(c.baseURL, method, endpoint, kp)
}

func (c *Connector) AccountSigned(kp keypair.KP, account string) (*Account, error) {
	endpoint := fmt.Sprintf("/accounts/%s", account)
	response, err := c.doSigned(kp, "GET", endpoint)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	switch response.StatusCode {
	case 404:
		return nil, nil
	case 200:
		var account Account
		err = json.NewDecoder(response.Body).Decode(&account)
		return &account, err
	default:
		return nil, fmt.Errorf("failed to load account: %d", response.StatusCode)
	}
}

func (c *Connector) BalanceAsset(balanceID string) (*Asset, error) {
	response, err := c.do("GET", fmt.Sprintf("/balances/%s/asset", balanceID))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 404:
		return nil, nil
	case 200:
		var asset Asset
		err = json.NewDecoder(response.Body).Decode(&asset)
		return &asset, err
	default:
		return nil, fmt.Errorf("failed to get balance asset: %d", response.StatusCode)
	}
}

func (c *Connector) Signers(accountID string) ([]Signer, error) {
	response, err := c.do("GET", fmt.Sprintf("/accounts/%s/signers", accountID))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 404:
		return nil, nil
	case 200:
		signers := struct {
			Signers []Signer `json:"signers"`
		}{}
		err = json.NewDecoder(response.Body).Decode(&signers)
		return signers.Signers, err
	default:
		return nil, fmt.Errorf("failed to load signers: %d", response.StatusCode)
	}
}

func (c *Connector) BalanceIDs(accountID string) (*BalanceIDResponse, error) {
	response, err := c.do("GET", fmt.Sprintf("/accounts/%s/balances", accountID))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 404:
		return nil, nil
	case 200:
		var result BalanceIDResponse
		err = json.NewDecoder(response.Body).Decode(&result)
		return &result, err
	default:
		return nil, fmt.Errorf("failed to load balances: %d", response.StatusCode)
	}
}

func (c *Connector) CoinEmissionRequests(kp keypair.KP, params *CoinEmissionRequestsParams) ([]CoinEmissionRequest, error) {
	endpoint := fmt.Sprintf("/coins_emission_requests?reference=%s&exchange=%s", params.Reference, params.Exchange)
	response, err := c.doSigned(kp, "GET", endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		cerResponse := CoinEmissionRequestsResponse{}

		err := json.NewDecoder(response.Body).Decode(&cerResponse)
		if err != nil {
			return nil, err
		}
		return cerResponse.Embedded.Records, nil
	case 404:
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to get cer: %d", response.StatusCode)
	}
}

func (c *Connector) ManageAssetPairOp(base string, quote string, physicalPrice int64) (xdr.Operation, error) {
	return xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageAssetPair,
			ManageAssetPairOp: &xdr.ManageAssetPairOp{
				Action:                  xdr.ManageAssetPairActionUpdatePrice,
				Base:                    xdr.AssetCode(base),
				Quote:                   xdr.AssetCode(quote),
				PhysicalPrice:           xdr.Int64(physicalPrice),
				PhysicalPriceCorrection: xdr.Int64(0),
				MaxPriceStep:            xdr.Int64(0),
				Policies:                xdr.Int32(0),
			},
		},
	}, nil
}

func (c *Connector) CreateBalanceOp(accountID, asset string) (xdr.Operation, error) {
	var op xdr.Operation

	var xAccountID xdr.AccountId
	err := xAccountID.SetAddress(accountID)
	if err != nil {
		return op, err
	}
	op = xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageBalance,
			ManageBalanceOp: &xdr.ManageBalanceOp{

				Action:      xdr.ManageBalanceActionCreate,
				Destination: xAccountID,
				Asset:       xdr.AssetCode(asset),
			},
		},
	}

	return op, nil
}

func (c *Connector) ReviewPaymentRequestOp(paymentID int64, accept bool) (xdr.Operation, error) {
	op := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeReviewPaymentRequest,
			ReviewPaymentRequestOp: &xdr.ReviewPaymentRequestOp{
				PaymentId: xdr.Uint64(paymentID),
				Accept:    accept,
			},
		},
	}

	return op, nil
}

func (c *Connector) submit(request *http.Request, tx string) ([]byte, error) {
	payload := resources.TxSubmission{
		TX: tx,
	}
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode payload")
	}

	// FIXME hazardous
	buf := b.Bytes()
	request.Body = ioutil.NopCloser(bytes.NewReader(buf))
	request.ContentLength = int64(b.Len())
	request.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(buf)
		return ioutil.NopCloser(r), nil
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body")
	}

	switch response.StatusCode {
	case 200:
		return body, nil
	default:
		serr, err := NewSubmitError(response.StatusCode, body)
		if err != nil {
			return body, err
		}
		return body, serr
	}
}

func (c *Connector) SubmitTX(tx string) error {
	request, err := NewRequest(c.baseURL, "POST", "/transactions")
	if err != nil {
		return errors.Wrap(err, "failed to build request")
	}
	_, err = c.submit(request, tx)
	if err != nil {
		return errors.Wrap(err, "failed to submit tx")
	}
	return nil
}

func (c *Connector) SubmitTXSignedVerbose(tx string, kp keypair.KP) ([]byte, error) {
	request, err := NewSignedRequest(c.baseURL, "POST", "/transactions", kp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}
	body, err := c.submit(request, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to submit tx")
	}
	return body, nil
}

type TransactionSuccess struct {
	Hash   string `json:"hash"`
	Ledger int32  `json:"ledger"`
	Env    string `json:"envelope_xdr"`
	Result string `json:"result_xdr"`
	Meta   string `json:"result_meta_xdr"`
}

func (c *Connector) SubmitTXVerbose(tx string) (*TransactionSuccess, error) {
	url, err := url.Parse(fmt.Sprintf("%s/transactions", c.baseURL))
	if err != nil {
		return nil, err
	}

	query := url.Query()
	query.Set("tx", tx)
	url.RawQuery = query.Encode()

	response, err := c.client.Post(url.String(), "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		var result TransactionSuccess
		err := json.NewDecoder(response.Body).Decode(&result)
		return &result, err
	default:
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read body")
		}
		serr, err := NewSubmitError(response.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, serr
	}
}

func (c *Connector) Transaction(tx *TransactionBuilder) *TransactionBuilder {
	if tx.Envelope != "" {
		env := xdr.TransactionEnvelope{}
		err := xdr.SafeUnmarshalBase64(tx.Envelope, &env)
		if err != nil {
			tx.err = err
			return tx
		}

		tx.Salt = uint64(env.Tx.Salt)

		tx.TimeBounds = &env.Tx.TimeBounds

		for _, signature := range env.Signatures {
			s, err := xdr.MarshalBase64(&signature)
			if err != nil {
				tx.err = err
				return tx
			}
			tx.Signatures = append(tx.Signatures, s)
		}
		tx.Source, err = keypair.Parse(env.Tx.SourceAccount.Address())
		if err != nil {
			tx.err = err
			return tx
		}

		tx.Operations = env.Tx.Operations
	}

	if tx.connector == nil {
		tx.connector = c
	}
	if tx.ops == nil {
		tx.ops = []Operation{}
	}
	return tx
}

func ParseBalanceID(addr string) (xdr.BalanceId, error) {
	raw, err := strkey.Decode(strkey.VersionByteBalanceID, addr)
	if err != nil {
		return xdr.BalanceId{}, err
	}

	if len(raw) != 32 {
		return xdr.BalanceId{}, errors.New("invalid address")
	}

	var ui xdr.Uint256
	copy(ui[:], raw)

	return xdr.NewBalanceId(xdr.CryptoKeyTypeKeyTypeEd25519, ui)
}
