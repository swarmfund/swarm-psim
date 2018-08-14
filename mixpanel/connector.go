package mixpanel

import (
	"net/url"
	"time"

	"github.com/dukex/mixpanel"
)

type Connector struct {
	mixpanelClient mixpanel.Mixpanel
}

func NewConnector(apiURL *url.URL, token string) *Connector {
	return &Connector{
		mixpanelClient: mixpanel.New(token, apiURL.String()),
	}
}

const mixpanelIPNotSpecified = "0"

// NewMixpanelEvent constructs event suitable to send to mixpanel using provided time
func newMixpanelEvent(time *time.Time) *mixpanel.Event {
	return &mixpanel.Event{
		IP:        mixpanelIPNotSpecified,
		Timestamp: time,
		Properties: map[string]interface{}{
			"Property": "Swarm Fund Invest",
		},
	}
}

func (c *Connector) SendEvent(account string, name string, time time.Time) error {
	return c.mixpanelClient.Track(account, name, newMixpanelEvent(&time))
}
