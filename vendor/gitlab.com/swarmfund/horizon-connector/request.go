package horizon

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/swarmfund/go/hash"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/xdr"
)

const (
	validUntilOffset        = 60
	EndpointForfeitRequests = "forfeit_requests"
)

// TODO come up with a better API for request builders

func NewRequest(baseURL, method, path string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", baseURL, path)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func NewSignedRequest(baseURL, method, path string, kp keypair.KP) (*http.Request, error) {
	req, err := NewRequest(baseURL, method, path)
	if err != nil {
		return nil, err
	}

	validUntil := fmt.Sprintf("%d", time.Now().Unix()+validUntilOffset)
	req.Header.Set("X-AuthValidUnTillTimestamp", validUntil)
	req.Header.Set("X-AuthPublicKey", kp.Address())

	base := fmt.Sprintf("{ uri: '%s', valid_untill: '%s'}", path, validUntil)
	hashBase := hash.Hash([]byte(base))

	decorated, err := kp.SignDecorated(hashBase[:])
	if err != nil {
		return nil, err
	}

	ohai, err := xdr.MarshalBase64(decorated)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-AuthSignature", ohai)

	return req, nil
}
