package internal

import (
	"context"
)

// Broadcaster is responsible for holding several targets and sending same events to them
type Broadcaster interface {
	AddTarget(target Target)
	BroadcastEvents(ctx context.Context, events <-chan ProcessedItem) chan error
}
