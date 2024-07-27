package resolvers

import (
	"context"
	"errors"
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
