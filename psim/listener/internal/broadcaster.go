package internal

import (
	"context"
)

type Source <-chan []BroadcastedEvent

// TODO review interface fields
type Broadcaster interface {
	SetSource(newSource Source)
	SetTargets(newTargets []Target)
	AddTarget(target Target)
	BroadcastEvents(ctx context.Context, events <-chan []BroadcastedEvent) error
}
