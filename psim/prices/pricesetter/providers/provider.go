package providers

import (
	"context"
	"sync"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
)

type provider struct {
	pendingPoints <-chan PricePoint
	name          string
	log           *logan.Entry

	pointsLock    *sync.Mutex
	points        []PricePoint
	lastSeenPoint PricePoint
}

func StartNewProvider(ctx context.Context, name string, pendingPoints <-chan PricePoint, log *logan.Entry) *provider {
	result := provider{
		pendingPoints: pendingPoints,
		name:          name,
		log:           log.WithField("price_provider", name),

		pointsLock:    new(sync.Mutex),
	}

	go result.fetchPointsInfinitely(ctx)

	return &result
}

// GetName - returns name of the provider
func (p *provider) GetName() string {
	return p.name
}

// GetPoints - returns points available
func (p *provider) GetPoints() []PricePoint {
	p.pointsLock.Lock()
	defer p.pointsLock.Unlock()

	result := make([]PricePoint, len(p.points))
	for i := range p.points {
		result[i] = p.points[i]
	}

	return result
}

// RemoveOldPoints - removes points which were created before minAllowedTime
func (p *provider) RemoveOldPoints(minAllowedTime time.Time) {
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

func (p *provider) fetchPointsInfinitely(ctx context.Context) {
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

func (p *provider) tryWriteNewPoint(point PricePoint) {
	// should be !Before to skip entries with same time
	if !p.lastSeenPoint.Time.Before(point.Time) {
		return
	}

	p.pointsLock.Lock()
	defer p.pointsLock.Unlock()

	p.points = append(p.points, point)
	p.lastSeenPoint = point
}
