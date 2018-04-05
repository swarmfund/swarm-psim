package horizon

import (
	"net/http"
	"net/url"

	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/account"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/asset"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/listener"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/operation"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/transaction"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/sale"
	"gitlab.com/tokend/keypair"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/user"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/balance"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/blob"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/document"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/templates"
)

type Connector struct {
	client *Client
}

func NewConnector(base *url.URL) *Connector {
	client := NewClient(http.DefaultClient, base)
	return &Connector{
		client,
	}
}

func (c *Connector) WithSigner(kp keypair.Full) *Connector {
	return &Connector{
		c.client.WithSigner(kp),
	}
}

func (c *Connector) Client() *Client {
	return c.client
}

func (c *Connector) Info() (info *Info, err error) {
	response, err := c.client.Get("/")
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if err := json.Unmarshal(response, &info); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal info")
	}
	return info, nil
}

func (c *Connector) Submitter() *Submitter {
	return &Submitter{
		client: c.client,
	}
}

func (c *Connector) Assets() *asset.Q {
	return asset.NewQ(c.client)
}

func (c *Connector) Accounts() *account.Q {
	return account.NewQ(c.client)
}

func (c *Connector) Transactions() *transaction.Q {
	return transaction.NewQ(c.client)
}

func (c *Connector) Sales() *sale.Q {
	return sale.NewQ(c.client)
}

func (c *Connector) Users() *user.Q {
	return user.NewQ(c.client)
}

func (c *Connector) Balances() *balance.Q {
	return balance.NewQ(c.client)
}

func (c *Connector) Listener() *listener.Q {
	// TODO Rename Operations to Requests? it does actually manages Requests only.
	return listener.NewQ(c.Transactions(), c.Operations())
}

// TODO Rename to Requests? it does actually manages Requests only.
func (c *Connector) Operations() *operation.Q {
	return operation.NewQ(c.client)
}

func (c *Connector) Blobs() *blob.Q {
	return blob.NewQ(c.client)
}

func (c *Connector) Documents() *documnet.Q {
	return documnet.NewQ(c.client)
}

func (c *Connector) Templates() *templates.Q{
	return templates.NewQ(c.client)
}
