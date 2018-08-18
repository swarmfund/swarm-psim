package eventsubmitter

import (
	"context"
	"sync"

	"gitlab.com/distributed_lab/logan/v3"
)

// BufferedTarget holds actual target and its events to broadcast
type BufferedTarget struct {
	Target Target
	Data   chan MaybeBroadcastedEvent
}

// GenericBroadcaster is a general-purpose Broadcaster implementation
type GenericBroadcaster struct {
	logger          *logan.Entry
	BufferedTargets []BufferedTarget
}

// NewGenericBroadcaster constructs a generic broadcaster with no targets
func NewGenericBroadcaster(logger *logan.Entry) *GenericBroadcaster {
	return &GenericBroadcaster{logger, []BufferedTarget{}}
}

const defaultTargetBufferSize = 1000

// AddTarget adds a target to broadcaster and initializes a channel for it
func (gb *GenericBroadcaster) AddTarget(target Target) {
	gb.BufferedTargets = append(gb.BufferedTargets, BufferedTarget{target, make(chan MaybeBroadcastedEvent, defaultTargetBufferSize)})
}

func (gb *GenericBroadcaster) putEventsToBufferedTargets(ctx context.Context, processedItems <-chan ProcessedItem) {
	targets := gb.BufferedTargets
	for _, target := range targets {
		target := target
		defer func() {
			close(target.Data)
		}()
	}
	for item := range processedItems {
		item := item
		for _, target := range targets {
			target := target
			select {
			case <-ctx.Done():
				return
			case target.Data <- MaybeBroadcastedEvent{BroadcastedEvent: item.BroadcastedEvent, Error: item.Error}:
				continue
			default:
				gb.logger.Warn("buffer busy, skiping")
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}
		}
	}
}

func (gb *GenericBroadcaster) sendEventsToBufferedTargets(ctx context.Context) {
	wg := new(sync.WaitGroup)
	for _, target := range gb.BufferedTargets {
		wg.Add(1)
		target := target
		go func(target BufferedTarget, ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					gb.logger.WithRecover(r).Error("panic while sending event to a target")
				}
				wg.Done()
			}()

			for event := range target.Data {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if event.Error != nil {
					gb.logger.WithError(event.Error).Warn("received invalid event")
					continue
				}

				err := target.Target.SendEvent(event.BroadcastedEvent)
				if err != nil {
					gb.logger.WithError(err).Error("failed to send event, skipping")
					continue
				}

				gb.logger.WithField("event_source", event.BroadcastedEvent.Account).WithField("event_name", event.BroadcastedEvent.Name).Info("sent event")
			}
		}(target, ctx)
	}
	wg.Wait()
}

// BroadcastEvents launches two goroutines - one copies events to buffered targets - second actually sends them to targets
func (gb *GenericBroadcaster) BroadcastEvents(ctx context.Context, items <-chan ProcessedItem) {
	go gb.putEventsToBufferedTargets(ctx, items)
	gb.sendEventsToBufferedTargets(ctx)
}
