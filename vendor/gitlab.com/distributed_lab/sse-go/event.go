package sse

import "io"

type Event struct {
	Type string
	Data io.Reader
	Err  error
}
