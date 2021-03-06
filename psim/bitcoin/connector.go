package bitcoin

import (
	"bytes"
	"net/http"

	"io/ioutil"
	"time"

	"encoding/json"

	"fmt"

	"strconv"

	"context"
	"encoding/base64"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient funds.")
	ErrAlreadyInChain    = errors.New("Transaction is already in chain.")
)

// Connector is interface Client uses to request some Bitcoin node, particularly Bitcoin Core.
type Connector interface {
	IsTestnet() bool

	// GetBlockCount must return index of last known Block
	GetBlockCount(context.Context) (uint64, error)
	// TODO Handle absent Block
	GetBlockHash(blockNumber uint64) (string, error)
	// GetBlock must return hex of Block
	// TODO Handle absent Block
	GetBlock(blockHash string) (string, error)

	GetBalance(includeWatchOnly bool) (float64, error)
	SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error)
	SendMany(addrToAmount map[string]float64) (resultTXHash string, err error)

	CreateRawTX(inputUTXOs []Out, outAddrToAmount map[string]float64) (resultTXHex string, err error)
	FundRawTX(initialTXHex, changeAddress string, includeWatching bool, feeRate *float64) (result *FundResult, err error)
	SignRawTX(initialTXHex string, inputUTXOs []InputUTXO, privateKeys []string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
	EstimateFee(blocks uint) (float64, error)
	GetTxUTXO(txHash string, vout uint32, unconfirmed bool) (*UTXO, error)
	ListUnspent(minConfirmations, maxConfirmations int, addresses []string) ([]WalletUTXO, error)
}

// NodeConnector is implementor of Connector interface,
// which requests Bitcoin core RPC to get the blockchain info
type NodeConnector struct {
	config  ConnectorConfig
	authKey string
	client  *http.Client
}

// NewNodeConnector returns new NodeConnector instance,
// created with provided ConnectorConfig.
func NewNodeConnector(config ConnectorConfig) Connector {
	return &NodeConnector{
		config:  config,
		authKey: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`%s:%s`, config.Node.User, config.Node.Password))),
		client: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
	}
}

// IsTestnet returns true if NodeConnector is using a testnet Bitcoin Node.
func (c *NodeConnector) IsTestnet() bool {
	return c.config.Testnet
}

// GetBlockCount returns index of a last known Block.
func (c *NodeConnector) GetBlockCount(ctx context.Context) (uint64, error) {
	var response struct {
		Response
		Result uint64 `json:"result"`
	}

	err := c.sendRequestWithCtx(ctx, "getblockcount", "", &response)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to send or parse get Block count request")
	}
	if response.Error != nil {
		return 0, errors.Wrap(response.Error, "Response for Block count request contains error")
	}

	return response.Result, nil
}

// GetBlockHash gets hash of Block by its index.
// TODO Handle absent Block
func (c *NodeConnector) GetBlockHash(blockNumber uint64) (string, error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	err := c.sendRequest("getblockhash", strconv.FormatUint(blockNumber, 10), &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse get Block hash request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for Block hash request contains error")
	}

	return response.Result, nil
}

// GetBlock gets raw Block by its hash and returns the raw Block encoded in hex.
// TODO Handle absent Block
func (c *NodeConnector) GetBlock(blockHash string) (string, error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	err := c.sendRequest("getblock", fmt.Sprintf(`"%s", false`, blockHash), &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse get Block request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for Block request contains error")
	}

	return response.Result, nil
}

func (c *NodeConnector) GetBalance(includeWatchOnly bool) (float64, error) {
	var response struct {
		Response
		Result float64 `json:"result"`
	}

	err := c.sendRequest("getbalance", fmt.Sprintf(`"", 1, %t`, includeWatchOnly), &response)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to send or parse get balance request")
	}
	if response.Error != nil {
		return 0, errors.Wrap(response.Error, "Response for get balance request contains error")
	}

	return response.Result, nil
}

func (c *NodeConnector) SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	// Empty strings in parameters stands for comments, `true` - is the subtract fee flag.
	err = c.sendRequest("sendtoaddress", fmt.Sprintf(`"%s", %.8f, "", "", true`, goalAddress, amount), &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse Send to Address request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for Send to Address request contains error")
	}

	return response.Result, nil
}

func (c *NodeConnector) SendMany(addrToAmount map[string]float64) (resultTXHash string, err error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	var lastAddr string
	params := `"", {` // FromAccount
	for addr, amount := range addrToAmount {
		lastAddr = addr
		params += fmt.Sprintf(`"%s": %.8f,`, addr, amount)
	}
	params = params[:len(params)-1] + fmt.Sprintf(`}, 1, "", ["%s"]`, lastAddr) // Confirmations, Comment, SubtractFeeFromAmount

	err = c.sendRequest("sendmany", params, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse SendMany request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for SendMany request contains error", logan.F{
			"params": params,
		})
	}

	return response.Result, nil
}

// CreateRawTX accepts nil as inputUTXOs.
// It is useful to pass empty inputUTXOs if you will later ask Wallet to Fund the TX (fill it with inputs).
func (c *NodeConnector) CreateRawTX(inputUTXOs []Out, outAddrToAmount map[string]float64) (resultTXHex string, err error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	// Inputs
	var inArrayParam string
	if len(inputUTXOs) == 0 {
		inArrayParam = `[]`
	} else {
		bb, err := json.Marshal(inputUTXOs)
		if err != nil {
			return "", errors.Wrap(err, "Failed to marshal inputUTXOs")
		}
		inArrayParam = string(bb)
	}

	// Outputs
	outsParam := `{`
	for addr, amount := range outAddrToAmount {
		outsParam += fmt.Sprintf(`"%s": %.8f,`, addr, amount)
	}
	outsParam = outsParam[:len(outsParam)-1] + `}`

	params := fmt.Sprintf(`%s, %s`, inArrayParam, outsParam)
	err = c.sendRequest("createrawtransaction", params, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse createRawTransaction request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for createRawTransaction request contains error")
	}

	return response.Result, nil
}

// FundRawTX runs fundrawtransaction request to the Bitcoin Node
// using flags `includeWatching` and `lockUnspents` as true.
// Fee is subtracted from the Output 0. Change position in Outputs is set to 1.
// If feeRate provided is nil - the wallet determines the fee.
//
// If Bitcoin Node returns -4:Insufficient funds error -
// ErrInsufficientFunds is returned.
//
// If returned error is nil - result is definitely not nil.
func (c *NodeConnector) FundRawTX(initialTXHex, changeAddress string, includeWatching bool, feeRate *float64) (result *FundResult, err error) {
	var response struct {
		Response
		Result FundResult `json:"result"`
	}

	params := fmt.Sprintf(`"%s", {
			"changeAddress": "%s",
			"changePosition": 1,
			"includeWatching": %t,
			"lockUnspents": true,
			"subtractFeeFromOutputs": [0]`, initialTXHex, changeAddress, includeWatching)
	if feeRate != nil {
		params = params + fmt.Sprintf(
			`,
		"feeRate": %.8f`, *feeRate)
	}
	params = params + `}`

	err = c.sendRequest("fundrawtransaction", params, &response)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send or parse fund raw Transaction request")
	}
	if response.Error != nil {
		if response.Error.Code == errCodeInsufficientFunds {
			return nil, ErrInsufficientFunds
		}

		return nil, errors.Wrap(response.Error, "Response for fund raw Transaction request contains error")
	}

	return &response.Result, nil
}

// SignRawTX signs the inputs of the provided TX with the provided privateKey.
// If the provided privateKey is nil - the TX will be tried to sign by Node, using
// the private keys Node owns.
func (c *NodeConnector) SignRawTX(initialTXHex string, outputsBeingSpent []InputUTXO, privateKeys []string) (resultTXHex string, err error) {
	var outsArray string
	if outputsBeingSpent == nil {
		outsArray = "[]"
	} else {
		bb, err := json.Marshal(outputsBeingSpent)
		if err != nil {
			return "", errors.Wrap(err, "Failed to marshal outputsBeingSpent")
		}

		outsArray = string(bb)
	}

	var response struct {
		Response
		Result struct {
			Hex      string `json:"hex"`
			Complete bool   `json:"complete"`
		} `json:"result"`
	}

	var privKeysParam string
	if len(privateKeys) == 0 {
		privKeysParam = `[]`
	} else {
		bb, err := json.Marshal(privateKeys)
		if err != nil {
			return "", errors.Wrap(err, "Failed to marshal private keys into JSON array")
		}
		privKeysParam = string(bb)
	}

	params := fmt.Sprintf(`"%s", %s, %s`, initialTXHex, outsArray, privKeysParam)
	err = c.sendRequest("signrawtransaction", params, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse signRawTransaction request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for sign raw Transaction request contains error")
	}

	return response.Result.Hex, nil
}

func (c *NodeConnector) SendRawTX(txHex string) (txHash string, err error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	err = c.sendRequest("sendrawtransaction", fmt.Sprintf(`"%s"`, txHex), &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse send raw Transaction request")
	}
	if response.Error != nil {
		fields := logan.F{
			"bitcoin_core_response_id":     response.ID,
			"bitcoin_core_response_result": response.Result,
		}

		if response.Error.Code == errCodeTransactionAlreadyInChain {
			return response.Result, errors.From(ErrAlreadyInChain, fields)
		}

		return "", errors.Wrap(response.Error, "Response for send raw Transaction request contains error", fields)
	}

	return response.Result, nil
}

func (c *NodeConnector) GetTxUTXO(txHash string, vout uint32, unconfirmed bool) (*UTXO, error) {
	var response struct {
		Response
		Result *UTXO `json:"result"`
	}

	err := c.sendRequest("gettxout", fmt.Sprintf(`"%s", %d, %t`, txHash, vout, unconfirmed), &response)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send or parse Get TX UTXO request")
	}
	if response.Error != nil {
		fields := logan.F{
			"bitcoin_core_response_id":     response.ID,
			"bitcoin_core_response_result": response.Result,
		}

		return nil, errors.Wrap(response.Error, "Response for Get TX UTXO request contains error", fields)
	}

	return response.Result, nil
}

func (c *NodeConnector) EstimateFee(blocks uint) (float64, error) {
	var response struct {
		Response
		Result float64 `json:"result"`
	}

	err := c.sendRequest("estimatefee", fmt.Sprintf(`%d`, blocks), &response)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to send or parse estimate fee request")
	}
	if response.Error != nil {
		return 0, errors.Wrap(response.Error, "Response for estimate fee request contains error")
	}

	return response.Result, nil
}

func (c *NodeConnector) ListUnspent(minConfirmations, maxConfirmations int, addresses []string) ([]WalletUTXO, error) {
	var response struct {
		Response
		Result []WalletUTXO `json:"result"`
	}

	var addressesString string
	if addresses == nil {
		addressesString = `[]`
	} else {
		bb, err := json.Marshal(addresses)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to marshal addresses into bytes")
		}
		addressesString = string(bb)
	}

	err := c.sendRequest("listunspent", fmt.Sprintf(`%d, %d, %s`, minConfirmations, maxConfirmations, addressesString), &response)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send or parse list unspent request")
	}
	if response.Error != nil {
		return nil, errors.Wrap(response.Error, "Response for list unspent request contains error")
	}

	return response.Result, nil
}

// DEPRECATED: use with ctx only
func (c *NodeConnector) sendRequest(methodName, params string, response interface{}) error {
	return c.sendRequestWithCtx(context.TODO(), methodName, params, response)
}

func (c *NodeConnector) sendRequestWithCtx(ctx context.Context, methodName, params string, response interface{}) error {
	// FIXME: stepko
	bodyStr := c.buildRequestBody("hardcoded_request_id", methodName, params)
	fields := logan.F{
		"node_url":     c.getNodeURL(),
		"request_body": bodyStr,
	}

	body := bytes.NewReader([]byte(bodyStr))
	request, err := http.NewRequest("POST", c.getNodeURL(), body)
	if err != nil {
		return errors.Wrap(err, "Failed to create http POST request", fields)
	}

	request.Header.Set("Authorization", "Basic "+c.authKey)
	if ctx != nil {
		request = request.WithContext(ctx)
	} else {
		request = request.WithContext(context.Background())
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "Failed to send request", fields)
	}
	fields["response_status_code"] = resp.StatusCode

	defer func() { _ = resp.Body.Close() }()
	respBodyBB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read response body", fields)
	}
	fields["raw_response_body"] = string(respBodyBB)

	err = json.Unmarshal(respBodyBB, response)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal response body to JSON", fields)
	}

	return nil
}

func (c *NodeConnector) getNodeURL() string {
	return fmt.Sprintf("http://%s:%d", c.config.Node.Host, c.config.Node.Port)
}

func (c *NodeConnector) buildRequestBody(requestID, methodName, params string) string {
	return `{"jsonrpc": "1.0", "id":"` + requestID + `", "method": "` + methodName + `", "params": [` + params + `] }`
}
