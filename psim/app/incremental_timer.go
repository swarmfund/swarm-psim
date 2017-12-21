package app

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type incrementalTimer struct {
	initialPeriod time.Duration
	multiplier    time.Duration

	currentPeriod time.Duration
	iteration     int
}

func newIncrementalTimer(initialPeriod time.Duration, multiplier int) *incrementalTimer {
	return &incrementalTimer{
		initialPeriod: initialPeriod,
		multiplier:    time.Duration(multiplier),

		currentPeriod: initialPeriod,
	}
}

func (t *incrementalTimer) next() <-chan time.Time {
	result := time.After(t.currentPeriod)

	t.currentPeriod = t.currentPeriod * t.multiplier
	t.iteration += 1

	return result
}

// TODO Comment
// TODO Add defer
//
// Runner function must do some job only once, iteration of job execution in loop is
// responsibility of RunOverIncrementalTimer.
//
// You are generally not supposed to log error inside runner, you should return error instead -
// errors returned from runner will be logged with stack.
//
// RunOverIncrementalTimer returns only returns if ctx is canceled.
// TODO Think of using different values for normalPeriod and initial abnormal period.
func RunOverIncrementalTimer(ctx context.Context, log *logan.Entry, runnerName string, runner func(context.Context) error, normalPeriod time.Duration) {
	log = log.WithField("runner", runnerName)
	normalTicker := time.NewTicker(normalPeriod)

	for {
		select {
		case <-ctx.Done():
			log.Info("Context is canceled - stopping runner.")
			return
		case <-normalTicker.C:
			if IsCanceled(ctx) {
				log.Info("Context is canceled - stopping runner.")
				return
			}

			err := runSafely(ctx, runner)

			if err != nil {
				log.WithStack(err).WithError(err).Error("Runner returned error.")

				runAbnormalExecution(ctx, log, runner, normalPeriod)
				if IsCanceled(ctx) {
					log.Info("Context is canceled - stopping runner.")
					return
				}
			}
		}
	}
}

// Only returns if runner returned nil error or ctx was canceled.
func runAbnormalExecution(ctx context.Context, log *logan.Entry, runner func(context.Context) error, initialPeriod time.Duration) {
	incrementalTimer := newIncrementalTimer(initialPeriod, 2)

	for {
		select {
		case <-ctx.Done():
			return
		case <-incrementalTimer.next():
			if IsCanceled(ctx) {
				return
			}

			err := runSafely(ctx, runner)
			if err == nil {
				log.Info("Runner is returning to normal execution.")
				return
			}
			log.WithField("retry_number", incrementalTimer.iteration).WithField("next_retry_period", incrementalTimer.currentPeriod).
				WithStack(err).WithError(err).Error("Runner returned error.")
		}
	}
}

func runSafely(ctx context.Context, runner func(context.Context) error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.Wrap(errors.FromPanic(rec), "Runner panicked")
		}
	}()

	return runner(ctx)
}
