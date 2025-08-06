package service

import (
	"context"

	private "github.com/egandro/news-deframer/gen/private"
)

// private service example implementation.
// The example methods log the requests and return zero values.
type privatesrvc struct{}

// NewPrivate returns the private service implementation.
func NewPrivate() private.Service {
	return &privatesrvc{}
}

// Ping implements ping.
func (s *privatesrvc) Ping(ctx context.Context) (res string, err error) {
	return "pong", nil
}
