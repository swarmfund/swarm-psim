package internal

import (
	"context"
	"fmt"
)

type GenericBroadcaster struct {
	Source      Source
	Targets     []Target
	TargetsData []chan BroadcastedEvent
}

func NewGenericBroadcaster(initialTarget Target) *GenericBroadcaster {
	return &GenericBroadcaster{nil, []Target{initialTarget}, make([]chan BroadcastedEvent, 1)}
}

func (b *GenericBroadcaster) SetSource(newSource Source) {
	b.Source = newSource
}

func (b *GenericBroadcaster) SetTargets(newTargets []Target) {
	b.Targets = newTargets
}

func (b *GenericBroadcaster) AddTarget(target Target) {
	b.TargetsData = append(b.TargetsData, make(chan BroadcastedEvent))
	b.Targets = append(b.Targets, target)
}

// TODO parallel targets
// TODO source refactoring
func (b *GenericBroadcaster) BroadcastEvents(ctx context.Context, eventsSource <-chan []BroadcastedEvent) error {
	for events := range eventsSource {
		for _, event := range events {
			for i := range b.TargetsData {
				go func() {
					b.TargetsData[i] <- event
				}()
			}
			for i := range b.Targets {
				fmt.Println(event)
				go func() {
					b.Targets[i].SendEvent(<-b.TargetsData[i])
				}()
			}
		}
	}
	return nil
}
