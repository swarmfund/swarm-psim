package internal

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// BufferedTarget holds actual target and its events to broadcast
type BufferedTarget struct {
	Target
	Data chan MaybeBroadcastedEvent
}

// GenericBroadcaster is a general-purpose Broadcaster implementation
type GenericBroadcaster struct {
	BufferedTargets []BufferedTarget
}

// NewGenericBroadcaster constructs a generic broadcaster with no targets
func NewGenericBroadcaster() *GenericBroadcaster {
	return &GenericBroadcaster{[]BufferedTarget{}}
}

// TODO target buffer size

// AddTarget adds a target to broadcaster and initializes a channel for it
func (b *GenericBroadcaster) AddTarget(target Target) {
	b.BufferedTargets = append(b.BufferedTargets, BufferedTarget{target, make(chan MaybeBroadcastedEvent, 1000)})
}

func putEventsToBufferedTargets(ctx context.Context, targets []BufferedTarget, processedItems <-chan ProcessedItem) {

	for item := range processedItems {
		item := item
		for _, target := range targets {
			target := target

			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case target.Data <- MaybeBroadcastedEvent{item.BroadcastedEvent, item.Error}:
				continue
			default:
				// TODO report err
			}
		}
	}
}

func sendEventsToBufferedTargets(ctx context.Context, targets []BufferedTarget) (errs chan error) {
	errs = make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(len(targets))

	go func() {
		wg.Wait()
		close(errs)
	}()

	for _, target := range targets {
		go func(target BufferedTarget) {
			defer wg.Done()

			for event := range target.Data {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if event.Error != nil {
					errs <- errors.Wrap(event.Error, "received invalid event")
					continue
				}

				err := target.SendEvent(event.BroadcastedEvent)
				if err != nil {
					errs <- errors.Wrap(err, "failed to send event to target, trying again")

					err := target.SendEvent(event.BroadcastedEvent)
					if err != nil {
						errs <- errors.Wrap(err, "second try failed, disabling target")
						return
					}
				}
			}
		}(target)
	}

	return errs
}

// BroadcastEvents launches two goroutines - one copies events to buffered targets - second actually sends them to targets
func (b *GenericBroadcaster) BroadcastEvents(ctx context.Context, items <-chan ProcessedItem) (errs chan error) {
	errs = make(chan error)

	go putEventsToBufferedTargets(ctx, b.BufferedTargets, items)

	go func() {
		defer close(errs)

		for err := range sendEventsToBufferedTargets(ctx, b.BufferedTargets) {
			if err != nil {
				errs <- errors.Wrap(err, "failed to send events to some targets")
			}
		}
	}()

	return errs
}
