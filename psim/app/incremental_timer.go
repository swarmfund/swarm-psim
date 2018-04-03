package app

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type incrementalTimer struct {
	initialPeriod time.Duration
	maxPeriod     time.Duration
	multiplier    time.Duration

	currentPeriod time.Duration
	iteration     int
}

func newIncrementalTimer(initialPeriod, maxPeriod time.Duration, multiplier int) *incrementalTimer {
	return &incrementalTimer{
		initialPeriod: initialPeriod,
		maxPeriod:     maxPeriod,
		multiplier:    time.Duration(multiplier),

		currentPeriod: initialPeriod,
		iteration:     0,
	}
}

func (t *incrementalTimer) next() <-chan time.Time {
	result := time.After(t.currentPeriod)

	t.currentPeriod = t.currentPeriod * t.multiplier

	// upper cap for timer
	if t.currentPeriod > t.maxPeriod {
		t.currentPeriod = t.maxPeriod
	}

	t.iteration += 1

	return result
}

// RunOverIncrementalTimer calls the runner with the normalPeriod, until the runner returns error.
// Once the runner returned error, it will be called with the abnormalPeriod,
// increasing the period in 2 times each retry.
// Once the runner returns nil(no error) in abnormal execution,
// it's execution comes back to the normal one and the runner
// is called with the normal Period again.
//
// If runner panics, the panic value will be converted to error and logged with stack.
//
// Runner function must do some job only once(not in a loop), iteration of job execution in loop is
// responsibility of RunOverIncrementalTimer func.
//
// You are generally not supposed to log error inside the runner,
// you should return error instead - errors returned from runner will be logged with stack.
//
// RunOverIncrementalTimer returns only if ctx is canceled.
// DEPRECATED, use package distributed_lab/running instead
func RunOverIncrementalTimer(ctx context.Context, log *logan.Entry, runnerName string, runner func(context.Context) error,
	normalPeriod time.Duration, abnormalPeriod time.Duration) {

	if normalPeriod == 0 {
		normalPeriod = 1
	}

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
				log.WithStack(err).WithError(err).Errorf("Runner '%s' returned error.", runnerName)

				runAbnormalExecution(ctx, log, runnerName, runner, abnormalPeriod)
				if IsCanceled(ctx) {
					log.Info("Context is canceled - stopping runner.")
					return
				}
			}
		}
	}
}

// RunUntilSuccess calls the runner again and again while the runner returns error.
// The time pause before the retry calling the runner is at first equal to initialRetryPeriod
// and becomes twice bigger each retry.
//
// If runner panics, the panic value will be converted to error and logged with stack.
//
// You are generally not supposed to log error inside the runner,
// you should return error instead - errors returned from runner will be logged with stack.
//
// RunUntilSuccess returns only if the runner returns nil or ctx canceled.
// TODO pass maxPausePeriod
// DEPRECATED, use package distributed_lab/running instead
func RunUntilSuccess(ctx context.Context, log *logan.Entry, runnerName string, runner func(context.Context) error, initialRetryPeriod time.Duration) {
	log = log.WithField("runner", runnerName)

	err := runner(ctx)
	if err == nil {
		// Brief success!
		return
	}

	incrementalTimer := newIncrementalTimer(initialRetryPeriod, 10*time.Minute, 2)

	for err != nil {
		log.WithField("retry_number", incrementalTimer.iteration).WithField("next_retry_after", incrementalTimer.currentPeriod).
			WithStack(err).WithError(err).Errorf("Runner '%s' returned error.", runnerName)

		select {
		case <-ctx.Done():
			return
		case <-incrementalTimer.next():
			if IsCanceled(ctx) {
				return
			}

			err = runSafely(ctx, runner)
		}
	}
}

// RunAbnormalExecution Only returns if runner returned nil error or ctx was canceled.
func runAbnormalExecution(ctx context.Context, log *logan.Entry, runnerName string, runner func(context.Context) error, initialRetryPeriod time.Duration) {
	incrementalTimer := newIncrementalTimer(initialRetryPeriod, 10*time.Minute, 2)

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
			log.WithField("retry_number", incrementalTimer.iteration).WithField("next_retry_after", incrementalTimer.currentPeriod).
				WithStack(err).WithError(err).Errorf("Runner '%s' returned error.", runnerName)
		}
	}
}

func runSafely(ctx context.Context, runner func(context.Context) error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.Wrap(errors.WithStack(errors.FromPanic(rec)), "Runner panicked")
		}
	}()

	return runner(ctx)
}
