package emails

import (
	"context"
	"sync"

	"gitlab.com/swarmfund/psim/psim/app"
)

type TaskSyncSet struct {
	mu   sync.Mutex
	data map[string]struct{}
}

func NewSyncSet() TaskSyncSet {
	return TaskSyncSet{
		mu:   sync.Mutex{},
		data: make(map[string]struct{}),
	}
}

func (s *TaskSyncSet) Put(ctx context.Context, new string) {
	put := func() <-chan struct{} {
		c := make(chan struct{})

		go func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.data[new] = struct{}{}

			close(c)
		}()

		return c
	}

	select {
	case <-ctx.Done():
		return
	case <-put():
		return
	}
}

func (s *TaskSyncSet) Exists(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	return ok
}

func (s *TaskSyncSet) Delete(values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, value := range values {
		delete(s.data, value)
	}
}

func (s *TaskSyncSet) Length() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.data)
}

func (s *TaskSyncSet) Range(ctx context.Context, f func(s string)) {
	// TODO Listen to ctx along with mutex
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.data {
		if app.IsCanceled(ctx) {
			return
		}

		f(key)
	}
}
