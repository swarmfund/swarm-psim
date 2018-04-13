package emails

import (
	"context"
	"sync"

	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/distributed_lab/notificator-server/client"
)

type Task struct {
	Destination string
	Subject     string
	Message     string
}

func (t Task) toPayload() notificator.EmailRequestPayload {
	return notificator.EmailRequestPayload {
		Destination: t.Destination,
		Subject: t.Subject,
		Message: t.Message,
	}
}

type TaskSyncSet struct {
	mu   sync.Mutex
	data map[Task]struct{}
}

func newSyncSet() TaskSyncSet {
	return TaskSyncSet{
		mu:   sync.Mutex{},
		data: make(map[Task]struct{}),
	}
}

func (s *TaskSyncSet) put(ctx context.Context, new Task) {
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

func (s *TaskSyncSet) delete(values []Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, value := range values {
		delete(s.data, value)
	}
}

func (s *TaskSyncSet) length() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.data)
}

func (s *TaskSyncSet) rangeThrough(ctx context.Context, f func(task Task)) {
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
