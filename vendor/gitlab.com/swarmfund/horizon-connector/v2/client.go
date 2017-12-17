package horizon

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/tokend/keypair"
)

type Client struct {
	base   *url.URL
	signer keypair.Full
}

func NewClient(base *url.URL) *Client {
	return &Client{
		base, nil,
	}
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	switch {
	case response.StatusCode == http.StatusUnauthorized:
		// TODO handle unathorized
		panic("not implemented")
	case response.StatusCode >= 500:
		// TODO handle server error
		panic("not implemented")
	default:
		// TODO only 2xx and 404 should fall here
		// pass down response for further processing
		return response, nil
	}
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}
	u = c.base.ResolveReference(u)
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}
	return c.Do(request)
}

func (c *Client) WithSigner(kp keypair.Full) *Client {
	return &Client{
		c.base, kp,
	}
}
