package providers

import (
	"context"
	"sync"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/prices/types"
)

type Provider struct {
	pendingPoints <-chan types.PricePoint
	name          string
	log           *logan.Entry

	pointsLock    *sync.Mutex
	points        []types.PricePoint
	lastSeenPoint types.PricePoint
}

func StartNewProvider(ctx context.Context, name string, pendingPoints <-chan types.PricePoint, log *logan.Entry) *Provider {
	result := Provider{
		pendingPoints: pendingPoints,
		name:          name,
		log:           log.WithField("price_provider", name),

		pointsLock: new(sync.Mutex),
	}

	go result.fetchPointsInfinitely(ctx)

	return &result
}

// GetName - returns name of the Provider
func (p *Provider) GetName() string {
	return p.name
}

// GetPoints - returns points available
func (p *Provider) GetPoints() []types.PricePoint {
	p.pointsLock.Lock()
	defer p.pointsLock.Unlock()

	result := make([]types.PricePoint, len(p.points))
	for i := range p.points {
		result[i] = p.points[i]
	}

	return result
}

// RemoveOldPoints - removes points which were created before minAllowedTime
func (p *Provider) RemoveOldPoints(minAllowedTime time.Time) {
	// it is guaranteed that prices added to slice in ascending order (entries with greater time comes last)
	// so we can just cut off part of the slice
	p.pointsLock.Lock()
	defer p.pointsLock.Unlock()
	newStartIndex := len(p.points)
	for i := range p.points {
		if p.points[i].Time.After(minAllowedTime) {
			newStartIndex = i
			break
		}
	}

	p.points = p.points[newStartIndex:len(p.points)]
}

func (p *Provider) fetchPointsInfinitely(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case point, ok := <-p.pendingPoints:
			if !ok {
				p.log.Warn("Pending Points chanel was closed - stopping.")
				return
			}

			p.tryWriteNewPoint(point)
		}
	}
}

func (p *Provider) tryWriteNewPoint(point types.PricePoint) {
	// should be !Before to skip entries with same time
	if !p.lastSeenPoint.Time.Before(point.Time) {
		return
	}

	p.pointsLock.Lock()
	defer p.pointsLock.Unlock()

	p.points = append(p.points, point)
	p.lastSeenPoint = point
}
