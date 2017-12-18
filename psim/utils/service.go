package utils

import "context"

type ServiceFn func(context.Context)

type Service interface {
	Run(ctx context.Context) chan error
}
