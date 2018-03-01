package airdrop

import (
	"context"
	"sync"

	"gitlab.com/swarmfund/psim/psim/app"
)

type SyncSet struct {
	mu   sync.Mutex
	data map[string]struct{}
}

func NewSyncSet() SyncSet {
	return SyncSet{
		mu:   sync.Mutex{},
		data: make(map[string]struct{}),
	}
}

func (s *SyncSet) Put(new string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[new] = struct{}{}
}

func (s *SyncSet) Exists(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	return ok
}

func (s *SyncSet) Delete(values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, value := range values {
		delete(s.data, value)
	}
}

func (s *SyncSet) Length() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.data)
}

func (s *SyncSet) Range(ctx context.Context, f func(s string)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.data {
		if app.IsCanceled(ctx) {
			return
		}

		f(key)
	}
}
