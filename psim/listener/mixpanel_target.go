package listener

import (
	"time"

	"github.com/dukex/mixpanel"
)

// DefaultMixpanelURL means "api.mixpanel.com"
const DefaultMixpanelURL = ""

// MixpanelTarget is used as Mixpanel client and target to broadcast events
type MixpanelTarget struct {
	mixpanel.Mixpanel
}

// NewMixpanelTarget constructs a mixpanel target and initializes a client
func NewMixpanelTarget(mixpanelToken string) *MixpanelTarget {
	return &MixpanelTarget{mixpanel.New(mixpanelToken, DefaultMixpanelURL)}
}

const mixpanelIPNotSpecified = "0"

// NewMixpanelEvent constructs event suitable to send to mixpanel using provided time
func NewMixpanelEvent(time *time.Time) *mixpanel.Event {
	return &mixpanel.Event{
		IP:        mixpanelIPNotSpecified,
		Timestamp: time,
		Properties: map[string]interface{}{
			"Property": "Swarm Fund Invest",
		},
	}
}

// SendEvent sends an event to mixpanel via mixpanel client
func (mt *MixpanelTarget) SendEvent(event BroadcastedEvent) error {
	return mt.Track(event.Account, string(event.Name), NewMixpanelEvent(&event.Time))
}
