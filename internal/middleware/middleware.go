package middleware

import (
	"context"
)

type HandlerFunc[T, R any] func(ctx context.Context, req T) (R, error)

type MiddlewareFunc[T, R any] func(next HandlerFunc[T, R]) HandlerFunc[T, R]

// Apply takes a handler and a variadic number of middleware and applies them in order.
func Apply[T, R any](h HandlerFunc[T, R], middlewares ...MiddlewareFunc[T, R]) HandlerFunc[T, R] {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
