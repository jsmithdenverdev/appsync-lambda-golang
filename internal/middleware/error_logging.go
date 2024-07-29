package middleware

import (
	"context"
	"errors"
	"log/slog"
)

func WithErrorLogging[T, R any](logger *slog.Logger, resolver string) MiddlewareFunc[T, R] {
	return func(next HandlerFunc[T, R]) HandlerFunc[T, R] {
		return func(ctx context.Context, req T) (R, error) {
			var r R
			r, err := next(ctx, req)
			if err != nil {
				logger.ErrorContext(ctx, "resolver failed", "resolver", resolver, "error", err)
				// TODO: Common error?
				return r, errors.New("internal server error")
			}
			return r, nil
		}
	}
}
