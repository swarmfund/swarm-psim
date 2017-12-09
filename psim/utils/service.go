package utils

import "context"

type ServiceFn func(context.Context)

type Service interface {
	Run() chan error
}
