package horizon

import (
	"net/http"
	"net/url"
	"time"

	"io/ioutil"

	"gitlab.com/swarmfund/horizon-connector/v2/internal/errors"
	"gitlab.com/tokend/keypair"
)

func throttle() chan time.Time {
	burst := 2 << 10
	ch := make(chan time.Time, burst)
	go func() {
		tick := time.Tick(3 * time.Second)
		// prefill buffer
		for i := 0; i < burst; i++ {
			ch <- time.Now()
		}
		for {
			select {
			case ch <- <-tick:
			default:
			}
		}
	}()
	return ch
}

type Client struct {
	base     *url.URL
	signer   keypair.Full
	throttle chan time.Time
	client   *http.Client
}

func NewClient(client *http.Client, base *url.URL) *Client {
	return &Client{
		base, nil, throttle(), client,
	}
}

func (c *Client) Do(request *http.Request) ([]byte, error) {
	<-c.throttle

	response, err := c.client.Do(request)
	if err != nil {
		return nil, errors.E(
			"failed to perform request",
			err,
			errors.Network,
			errors.Path(request.URL.String()),
		)
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.E(
			"failed to read response body",
			err,
			errors.Runtime,
			errors.Path(request.URL.String()),
		)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return bytes, nil
	case http.StatusNotFound:
		return nil, nil
	case http.StatusTooManyRequests:
		// TODO look at x-rate-limit headers and slow down
		panic("not implemented")
	case http.StatusBadRequest:
		return nil, errors.E(
			"request was invalid in some way",
			errors.Runtime,
			errors.Response(bytes),
			errors.Status(response.StatusCode),
			errors.Path(request.URL.String()),
		)
	case http.StatusUnauthorized:
		return nil, errors.E(
			"signer is not allowed to access resource",
			errors.Runtime,
			errors.Response(bytes),
			errors.Status(response.StatusCode),
			errors.Path(request.URL.String()),
		)
	default:
		return nil, errors.E(
			"something bad happened",
			errors.Runtime,
			errors.Response(bytes),
			errors.Status(response.StatusCode),
			errors.Path(request.URL.String()),
		)
	}
}

func (c *Client) Get(endpoint string) ([]byte, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.E(
			"failed to parse endpoint",
			err,
			errors.Runtime,
		)
	}

	u = c.base.ResolveReference(u)
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.E(
			"failed to build request",
			err,
			errors.Runtime,
		)
	}
	return c.Do(request)
}

func (c *Client) WithSigner(kp keypair.Full) *Client {
	return &Client{
		c.base, kp, c.throttle, http.DefaultClient,
	}
}
