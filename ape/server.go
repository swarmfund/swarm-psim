package ape

import (
	"context"
	"net"
	"net/http"
)

type Response interface{}
type PipelineFn func() Response
type Handler func(r *http.Request) Response

type HTTPServer struct {
	Middlewares []int
}

func (s *HTTPServer) ListenAndServe(ctx context.Context, listener net.Listener) {}

func Do(fns ...PipelineFn) Response {
	for _, fn := range fns {
		response := fn()
		if response != nil {
			return response
		}
	}
	return nil
}

func A(handler func(r *Request) Response) http.HandlerFunc {
	return nil
}

type Request struct{}
