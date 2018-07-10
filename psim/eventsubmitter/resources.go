package eventsubmitter

import (
	"gitlab.com/swarmfund/psim/psim/eventsubmitter/internal"
)

// BroadcastedEvent exported to be used in listener package
type BroadcastedEvent = internal.BroadcastedEvent

// MaybeBroadcastedEvent exported to be used in listener package
type MaybeBroadcastedEvent = internal.MaybeBroadcastedEvent

// BroadcastedEventName exported to be used in listener package
type BroadcastedEventName = internal.BroadcastedEventName

// Processor exported to be used in listener package
type Processor = internal.Processor

// OpData exported to be used in listener package
type OpData = internal.OpData

// Broadcaster exported to be used in listener package
type Broadcaster = internal.Broadcaster

// Target exported to be used in listener package
type Target = internal.Target

// Extractor exported to be used in listener package
type Extractor = internal.Extractor

// ExtractedItem exported to be used in listener package
type ExtractedItem = internal.ExtractedItem

// ProcessedItem exported to be used in listener package
type ProcessedItem = internal.ProcessedItem

// Handler exported to be used in listener package
type Handler = internal.Handler
