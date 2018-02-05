package bitcoin

import (
	"bytes"
	"net/http"

	"io/ioutil"
	"time"

	"encoding/json"

	"fmt"

	"strconv"

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
	GetBlockCount() (uint64, error)
	// TODO Handle absent Block
	GetBlockHash(blockIndex uint64) (string, error)
	// GetBlock must return hex of Block
	// TODO Handle absent Block
	GetBlock(blockHash string) (string, error)
	GetBalance(includeWatchOnly bool) (float64, error)
	SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error)
	SendMany(addrToAmount map[string]float64) (resultTXHash string, err error)
	CreateRawTX(addrToAmount map[string]float64) (resultTXHex string, err error)
	FundRawTX(initialTXHex, changeAddress string, includeWatching bool, feeRate *float64) (result *FundResult, err error)
	SignRawTX(initialTXHex string, inputUTXOs []Out, privateKey *string) (resultTXHex string, err error)
	SendRawTX(txHex string) (txHash string, err error)
}

type Out struct {
	TXHash       string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript *string `json:"redeemScript,omitempty"`
}

// NodeConnector is implementor of Connector interface,
// which requests Bitcoin core RPC to get the blockchain info
type NodeConnector struct {
	config ConnectorConfig
	client *http.Client
}

// NewNodeConnector returns new NodeConnector instance,
// created with provided ConnectorConfig.
func NewNodeConnector(config ConnectorConfig) Connector {
	return &NodeConnector{
		config: config,
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
func (c *NodeConnector) GetBlockCount() (uint64, error) {
	var response struct {
		Response
		Result uint64 `json:"result"`
	}

	err := c.sendRequest("getblockcount", "", &response)
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
func (c *NodeConnector) GetBlockHash(blockIndex uint64) (string, error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	err := c.sendRequest("getblockhash", strconv.FormatUint(blockIndex, 10), &response)
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

func (c *NodeConnector) CreateRawTX(addrToAmount map[string]float64) (resultTXHex string, err error) {
	var response struct {
		Response
		Result string `json:"result"`
	}

	params := `[], {` // Inputs
	for addr, amount := range addrToAmount {
		params += fmt.Sprintf(`"%s": %.8f,`, addr, amount)
	}
	params = params[:len(params)-1] + `}`

	err = c.sendRequest("createrawtransaction", params, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse createRawTransaction request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for createRawTransaction request contains error")
	}

	return response.Result, nil
}

type FundResult struct {
	Hex            string  `json:"hex"`
	ChangePosition int     `json:"changepos"`
	FeePaid        float64 `json:"fee"`
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
func (c *NodeConnector) SignRawTX(initialTXHex string, outputsBeingSpent []Out, privateKey *string) (resultTXHex string, err error) {
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

	privateKeys := `[`
	if privateKey != nil {
		privateKeys = privateKeys + fmt.Sprintf(`"%s"`, *privateKey)
	}
	privateKeys = privateKeys + `]`

	params := fmt.Sprintf(`"%s", %s, %s`, initialTXHex, outsArray, privateKeys)

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

func (c *NodeConnector) sendRequest(methodName, params string, response interface{}) error {
	request, err := c.buildRequest("hardcoded_request_id", methodName, params)
	if err != nil {
		return errors.Wrap(err, "Failed to build request")
	}

	request.Header.Set("Authorization", "Basic "+c.config.NodeAuthKey)

	resp, err := c.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "Failed to send request")
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read response body")
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal response body to JSON")
	}

	return nil
}

func (c *NodeConnector) buildRequest(requestID, methodName, params string) (*http.Request, error) {
	bodyStr := c.buildRequestBody(requestID, methodName, params)
	body := bytes.NewReader([]byte(bodyStr))

	request, err := http.NewRequest("POST", c.getNodeURL(), body)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (c *NodeConnector) getNodeURL() string {
	return fmt.Sprintf("http://%s:%d", c.config.NodeIP, c.config.NodePort)
}

func (c *NodeConnector) buildRequestBody(requestID, methodName, params string) string {
	return `{"jsonrpc": "1.0", "id":"` + requestID + `", "method": "` + methodName + `", "params": [` + params + `] }`
}
