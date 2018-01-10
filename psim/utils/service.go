package utils

import "context"

// TODO Move to app
type Service interface {
	// TODO Stop returning chan error, do blocking Run() methods instead
	Run(ctx context.Context) chan error
}
