package handlers

import (
	"context"
	"log/slog"

	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/requests"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/responses"
)

func HandleCreateItem(
	logger *slog.Logger,
) func(ctx context.Context, request requests.CreateItem) (responses.CreateItem, error) {
	return func(ctx context.Context, request requests.CreateItem) (responses.CreateItem, error) {
		logger.InfoContext(ctx, "create item", "request", request)

		return responses.CreateItem{
			Item: responses.Item{
				Name: request.Name,
			},
		}, nil
	}
}
