package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type BatchInvokeResponse[T any] struct {
	Value T     `json:"value"`
	Error error `json:"error,omitempty"`
}

func WithBatchInvoke[T, R any](next func(ctx context.Context, req []R) ([]T, []error)) func(ctx context.Context, req []R) ([]BatchInvokeResponse[T], error) {
	return func(ctx context.Context, req []R) ([]BatchInvokeResponse[T], error) {
		results := make([]BatchInvokeResponse[T], len(req))

		nextResults, errs := next(ctx, req)

		if len(nextResults) != len(req) || len(errs) != len(req) {
			return results, errors.New("the batch function supplied did not return an array of responses the same length as the array of keys")
		}

		for i := 0; i < len(req); i++ {
			results[i].Value = nextResults[i]
			results[i].Error = errs[i]
		}

		return results, nil
	}
}

func resolveError[T any](name string, logger *slog.Logger) func(ctx context.Context, err error) (T, error) {
	return func(ctx context.Context, err error) (T, error) {
		logger.ErrorContext(ctx, fmt.Sprintf("resolver %s failed", name), "error", err)
		var t T
		return t, errInternalServer
	}
}
