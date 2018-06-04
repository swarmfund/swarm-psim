package listener

import (
	"time"

	"github.com/dukex/mixpanel"
)

const defaultMixpanelURL = ""

// MixpanelTarget is used as Mixpanel client and target to broadcast events
type MixpanelTarget struct {
	mixpanel.Mixpanel
}

func NewMixpanelTarget(mixpanelToken string) *MixpanelTarget {
	return &MixpanelTarget{mixpanel.New(mixpanelToken, defaultMixpanelURL)}
}

// TODO utilize IP?
const mixpanelIpNotSpecified = "0"

// DEPRECATED - Events should get their time from Tx
func mixpanelCurrentTime() *time.Time {
	return nil
}

// TODO mixpanel refuses using time
func NewMixpanelEvent(time *time.Time) *mixpanel.Event {
	return &mixpanel.Event{
		IP:        mixpanelIpNotSpecified,
		Timestamp: time,
		Properties: map[string]interface{}{
			"StubProperty": "StubValue",
		},
	}
}

func (mt *MixpanelTarget) SendEvent(event BroadcastedEvent) error {
	return mt.Track(event.Account, string(event.Name), NewMixpanelEvent(event.Time))
}
