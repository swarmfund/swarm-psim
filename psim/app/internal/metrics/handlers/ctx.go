package handlers

import (
	"context"
	"net/http"
	"sync"

	"gitlab.com/swarmfund/psim/psim/app/internal/data"
)

type ctxKey int

const (
	mutexCtxKey ctxKey = iota
	serviceCtxKey
)

func CtxMutex(mutex *sync.Mutex) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, mutexCtxKey, mutex)
	}
}

func Mutex(r *http.Request) *sync.Mutex {
	return r.Context().Value(mutexCtxKey).(*sync.Mutex)
}

func CtxServices(services []data.MetricsKeeper) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, serviceCtxKey, services)
	}
}

func Services(r *http.Request) []data.MetricsKeeper {
	return r.Context().Value(serviceCtxKey).([]data.MetricsKeeper)
}
