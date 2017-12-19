package bitcoin

import (
	"bytes"
	"net/http"

	"io/ioutil"
	"time"

	"encoding/json"

	"fmt"

	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Connector is interface Client uses to request some Bitcoin node, particularly Bitcoin Core.
type Connector interface {
	IsTestnet() bool
	// GetBlockCount must return index of last known Block
	GetBlockCount() (uint64, error)
	GetBlockHash(blockIndex uint64) (string, error)
	// GetBlock must return hex of Block
	GetBlock(blockHash string) (string, error)
	GetBalance() (float64, error)
	SendToAddress(goalAddress string, amount float64) (resultTXHash string, err error)
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

func (c *NodeConnector) GetBalance() (float64, error) {
	var response struct {
		Response
		Result float64 `json:"result"`
	}

	err := c.sendRequest("getbalance", "", &response)
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
	err = c.sendRequest("sendtoaddress", fmt.Sprintf(`"%s", %f, "", "", true`, goalAddress, amount), &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to send or parse send to Address request")
	}
	if response.Error != nil {
		return "", errors.Wrap(response.Error, "Response for send to Address request contains error")
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
