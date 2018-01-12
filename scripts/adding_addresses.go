package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Response struct {
	ID    string `json:"id"`
	Error *Error `json:"error"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Code, e.Message)
}

func main() {
	log := logan.New()

	args := os.Args[1:]
	if len(args) < 4 {
		log.Panic("Need Node url(1), auth key(2) and extPrivKey(3), n(4) to be passed as command line arguments.")
	}
	url := args[0]
	authKey := args[1]
	extPrivKey := args[2]

	n, err := strconv.Atoi(args[3])
	if err != nil {
		log.WithError(err).Panic("Failed to parse integer from the fourth argument")
	}

	params := &chaincfg.TestNet3Params
	privKeys, err := derivePrivateKeys(extPrivKey, params, n)
	if err != nil {
		log.WithError(err).Panic("Failed to derive private keys from extended key.")
		return
	}

	err = importPrivateKeys(log, url, authKey, privKeys)
	if err != nil {
		log.WithError(err).Panic("Failed to import PrivateKeys.")
	}
}

func derivePrivateKeys(extPrivKey string, params *chaincfg.Params, n int) ([]string, error) {
	extendedKey, err := hdkeychain.NewKeyFromString(extPrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create extended Key from base58 extended private key")
	}

	var result []string
	for i := 0; i < n; i++ {
		childKey, err := extendedKey.Child(hdkeychain.HardenedKeyStart + uint32(i))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to derive child key", logan.F{"i": i})
		}

		privKey, err := childKey.ECPrivKey()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get ECPrivKey of the derived Child Key", logan.F{"i": i})
		}

		//result = append(result, toWalletImportFormat(privKey, params, true))
		wif, err := btcutil.NewWIF(privKey, params, true)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create new WIF from private key", logan.F{"i": i})
		}

		result = append(result, wif.String())
	}

	return result, nil
}

func importPrivateKeys(log *logan.Entry, url string, authKey string, privateKeys []string) error {
	for i, privKey := range privateKeys {
		if privKey == "" {
			continue
		}

		err := sendRequestToBTCNode(url, authKey, "importprivkey", fmt.Sprintf(`"%s", "", false`, privKey))
		if err != nil {
			return errors.Wrap(err, "Failed to import private key", logan.F{"i": i})
		}

		log.WithField("i", i).Debug("Imported private key successfully.")
	}

	return nil
}

func sendRequestToBTCNode(url, authKey, methodName, params string) error {
	request, err := buildRequest(url, "hardcoded_request_id", methodName, params)
	if err != nil {
		return errors.Wrap(err, "Failed to build request")
	}

	request.Header.Set("Authorization", "Basic "+authKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "Failed to send request")
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read response body")
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal response body to JSON", logan.F{
			"status_code":       resp.StatusCode,
			"raw_response_body": string(body),
		})
	}

	if response.Error != nil {
		return errors.Wrap(err, "Node returned non nil error", logan.F{
			"status_code": resp.StatusCode,
		})
	}

	return nil
}

func buildRequest(url, requestID, methodName, params string) (*http.Request, error) {
	bodyStr := buildRequestBody(requestID, methodName, params)
	body := bytes.NewReader([]byte(bodyStr))

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func buildRequestBody(requestID, methodName, params string) string {
	return `{"jsonrpc": "1.0", "id":"` + requestID + `", "method": "` + methodName + `", "params": [` + params + `] }`
}
