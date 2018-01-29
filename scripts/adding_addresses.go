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

	"github.com/btcsuite/btcd/btcec"
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
	if len(args) < 5 {
		log.Panic("Need Bitcoin core url(1), Bitcoin core auth key(2), extPrivKey(3), index of first private key to import(4) and index of last private key to import(5) to be passed as command line arguments.\n" +
			"Example: http://127.0.0.1:8332 base64({{username}}:{{password}}) xprv....... 0 10000")
	}
	url := args[0]
	authKey := args[1]
	extPrivKey := args[2]

	var params *chaincfg.Params
	switch extPrivKey[:4] {
	case "xprv":
		params = &chaincfg.MainNetParams
	case "tprv":
		params = &chaincfg.TestNet3Params
	default:
		panic("Private key starts with neither 'xprv'(main net) not 'tprv'(test net).")
	}

	firstIndex, err := strconv.Atoi(args[3])
	if err != nil {
		log.WithError(err).Panic("Failed to parse integer from the argument (4) (index of first private key to import).")
	}

	lastIndex, err := strconv.Atoi(args[4])
	if err != nil {
		log.WithError(err).Panic("Failed to parse integer from the argument (5) (index of last private key to import).")
	}

	privKeys, err := derivePrivateKeys(extPrivKey, firstIndex, lastIndex, false)
	if err != nil {
		log.WithError(err).Panic("Failed to derive private keys from extended key.")
		return
	}

	log.Info("Finished deriving private keys from extended private key.")

	err = importPrivateKeys(log, url, params, authKey, privKeys, firstIndex)
	if err != nil {
		log.WithError(err).Error("Failed to import PrivateKeys.")
		return
	}
}

func derivePrivateKeys(extPrivKey string, firstIndex, lastIndex int, hardened bool) ([]*btcec.PrivateKey, error) {
	extendedKey, err := hdkeychain.NewKeyFromString(extPrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create extended Key from base58 extended private key")
	}

	var result []*btcec.PrivateKey
	for i := firstIndex; i < lastIndex+1; i++ {
		childID := uint32(i)
		if hardened {
			childID += hdkeychain.HardenedKeyStart
		}

		childKey, err := extendedKey.Child(childID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to derive child key", logan.F{"i": i})
		}

		privKey, err := childKey.ECPrivKey()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get ECPrivKey of the derived Child Key", logan.F{"i": i})
		}

		result = append(result, privKey)
	}

	return result, nil
}

func importPrivateKeys(log *logan.Entry, url string, params *chaincfg.Params, authKey string, privateKeys []*btcec.PrivateKey, firstIndex int) error {
	for i, privKey := range privateKeys {
		if privKey == nil {
			continue
		}

		var addr string
		addrPKH, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey.PubKey().SerializeCompressed()), params)
		if err != nil {
			addr = err.Error()
		} else {
			if addrPKH != nil {
				addr = addrPKH.String()
			}
		}

		fields := logan.F{
			"i":    firstIndex + i,
			"addr": addr,
		}

		wif, err := btcutil.NewWIF(privKey, params, true)
		if err != nil {
			return errors.Wrap(err, "Failed to create new WIF from private key", fields)
		}

		err = sendRequestToBTCNode(url, authKey, "importprivkey", fmt.Sprintf(`"%s", "", false`, wif.String()))
		for err != nil {
			// Retry until success
			log.WithFields(fields).WithError(err).Warn("Failed to import private key, will now sleep 10s and retry.")
			time.Sleep(10 * time.Second)

			err = sendRequestToBTCNode(url, authKey, "importprivkey", fmt.Sprintf(`"%s", "", false`, wif.String()))
		}

		log.WithFields(fields).Debug("Imported private key successfully.")
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
