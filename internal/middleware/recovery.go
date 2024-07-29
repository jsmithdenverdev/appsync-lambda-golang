package middleware

import (
	"context"
	"log/slog"
)

func WithRecovery[T, R any](logger *slog.Logger, resolver string) MiddlewareFunc[T, R] {
	return func(next HandlerFunc[T, R]) HandlerFunc[T, R] {
		return func(ctx context.Context, req T) (R, error) {
			defer func() {
				if err := recover(); err != nil {
					logger.ErrorContext(ctx, "failed to recover from resolver panic", "resolver", resolver, "error", err)
				}
			}()
			return next(ctx, req)
		}
	}
}
