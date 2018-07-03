package reporter

import (
	"gitlab.com/swarmfund/psim/psim/listener/internal"
)

// BroadcastedEvent exported to be used in listener package
type BroadcastedEvent = internal.BroadcastedEvent

// MaybeBroadcastedEvent exported to be used in listener package
type MaybeBroadcastedEvent = internal.MaybeBroadcastedEvent

// BroadcastedEventName exported to be used in listener package
type BroadcastedEventName = internal.BroadcastedEventName

// Broadcaster exported to be used in listener package
type Broadcaster = internal.Broadcaster

// Target exported to be used in listener package
type Target = internal.Target

// ProcessedItem exported to be used in listener package
type ProcessedItem = internal.ProcessedItem
