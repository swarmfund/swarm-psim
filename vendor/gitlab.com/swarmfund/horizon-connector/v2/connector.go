package horizon

import (
	"net/http"
	"net/url"

	"gitlab.com/swarmfund/horizon-connector/v2/internal/account"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/asset"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/listener"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/operation"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/transaction"
	"gitlab.com/tokend/keypair"
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

func (c *Connector) Assets() *asset.Q {
	return asset.NewQ(c.client)
}

func (c *Connector) Accounts() *account.Q {
	return account.NewQ(c.client)
}

func (c *Connector) Transactions() *transaction.Q {
	return transaction.NewQ(c.client)
}

func (c *Connector) Listener() *listener.Q {
	// TODO Rename Operations to Requests? it does actually manages Requests only.
	return listener.NewQ(c.Transactions(), c.Operations())
}

// TODO Rename to Requests? it does actually manages Requests only.
func (c *Connector) Operations() *operation.Q {
	return operation.NewQ(c.client)
}
