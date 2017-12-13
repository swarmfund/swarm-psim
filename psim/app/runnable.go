package app

import (
	"sync"
	"context"
)

// Runner is used in StartRunners.
type Runner interface {
	// Run should be a blocking method
	// and return only when `Runner` and all of its child are finished.
	// Run must stop if passed ctx is canceled.
	Run(ctx context.Context)
}

// StartRunners runs each provided runner is a separate goroutine,
// incrementing the provided finishWaiter on 1 before each runner start.
//
// StartRunners is not blocking.
//
// Once a runner's Run() method finishes - the finishWaiter is decreased by 1.
// So finishWaiter will be fully done, when all runners are done.
func StartRunners(ctx context.Context, finishWaiter *sync.WaitGroup, runners ...Runner) {
	for r := range runners {
		finishWaiter.Add(1)

		// Yes this copying is needed, because all the goroutines will be started after this loop finishes.
		ohigo := r
		go func() {
			defer finishWaiter.Done()
			runners[ohigo].Run(ctx)
		}()
	}
}
