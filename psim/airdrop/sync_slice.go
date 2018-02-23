package airdrop

import (
	"sync"
	"context"
	"gitlab.com/swarmfund/psim/psim/app"
)

type SyncSet struct {
	mu   sync.Mutex
	data map[string]struct{}
}

func (s *SyncSet) Put(new string) {
	s.mu.Lock()
	defer func() { s.mu.Unlock() }()

	s.data[new] = struct{}{}
}

func (s *SyncSet) Range(ctx context.Context, f func(s string)) {
	s.mu.Lock()
	defer func() { s.mu.Unlock() }()

	for key := range s.data {
		if app.IsCanceled(ctx) {
			return
		}

		f(key)
	}
}

func (s *SyncSet) Delete(values []string) {
	s.mu.Lock()
	defer func() { s.mu.Unlock() }()

	for _, value := range values {
		delete(s.data, value)
	}
}
