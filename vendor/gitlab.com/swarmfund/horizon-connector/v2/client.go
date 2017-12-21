package horizon

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
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
}

func NewClient(base *url.URL) *Client {
	return &Client{
		base, nil, throttle(),
	}
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	<-c.throttle
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
		c.base, kp, c.throttle,
	}
}
