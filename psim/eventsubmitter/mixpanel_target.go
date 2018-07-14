package eventsubmitter

import (
	"gitlab.com/swarmfund/psim/mixpanel"
)

// DefaultMixpanelURL means "api.mixpanel.com"
const DefaultMixpanelURL = ""

// MixpanelTarget is used as Mixpanel client and target to broadcast events
type MixpanelTarget struct {
	*mixpanel.Connector
}

// NewMixpanelTarget constructs a mixpanel target and initializes a client
func NewMixpanelTarget(connector *mixpanel.Connector) *MixpanelTarget {
	return &MixpanelTarget{connector}
}

// SendEvent sends an event to mixpanel via mixpanel client
func (mt *MixpanelTarget) SendEvent(event *BroadcastedEvent) error {
	return mt.Connector.SendEvent(event.Account, string(event.Name), event.Time)
}
