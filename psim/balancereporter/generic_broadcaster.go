package balancereporter

import (
	"context"
	"sync"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/regources"
)

type Target interface {
	SendEvent(event *regources.BalancesReport, swmAmount int64, threshold int64, date time.Time) (err error)
}

// BufferedTarget holds actual target and its events to broadcast
type BufferedTarget struct {
	Target Target
	Data   chan BroadcastedReport
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
	gb.BufferedTargets = append(gb.BufferedTargets, BufferedTarget{target, make(chan BroadcastedReport, defaultTargetBufferSize)})
}

func (gb *GenericBroadcaster) putEventsToBufferedTargets(ctx context.Context, processedItems <-chan BroadcastedReport) {
	for item := range processedItems {
		item := item
		for _, target := range gb.BufferedTargets {
			target := target
			select {
			case <-ctx.Done():
				return
			case target.Data <- item:
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

				err := target.Target.SendEvent(event.Report, event.SWMAmount, event.Threshold, event.Date)
				if err != nil {
					gb.logger.WithError(err).Error("failed to send event, skipping")
					continue
				}

			}
		}(target, ctx)
	}
	wg.Wait()
}

type BroadcastedReport struct {
	Report    *regources.BalancesReport
	SWMAmount int64
	Threshold int64
	Date      time.Time
}

// BroadcastEvents launches two goroutines - one copies events to buffered targets - second actually sends them to targets
func (gb *GenericBroadcaster) BroadcastEvents(ctx context.Context, items <-chan BroadcastedReport) {
	go gb.putEventsToBufferedTargets(ctx, items)
	gb.sendEventsToBufferedTargets(ctx)
}
