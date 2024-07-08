package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	lambda.Start(handler(logger))
}

func handler(logger *slog.Logger) func(ctx context.Context, event interface{}) (interface{}, error) {
	return func(ctx context.Context, event interface{}) (interface{}, error) {
		logger.InfoContext(ctx, "handler called", "event", event)
		return struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}{
			ID:          "1",
			Name:        "Example",
			Description: "An example",
		}, nil
	}
}
