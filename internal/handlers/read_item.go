package handlers

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/responses"
)

// HandleReadItemConfig represents the config needed by the ReadItem Lambda
// function.
type HandleReadItemConfig struct {
	TableName string
}

// handleReadItemDynamoDBClient represents the Dynamo DB operations that can
// be performed from HandleReadItem.
type handleReadItemDynamoDBClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// HandleReadItem is the handler for the ReadItem Lambda function.
func HandleReadItem(
	config HandleReadItemConfig,
	logger *slog.Logger,
	dynamodbclient handleReadItemDynamoDBClient,
) func(ctx context.Context, ids []string) ([]responses.Item, error) {
	return func(ctx context.Context, ids []string) ([]responses.Item, error) {
		var response []responses.Item

		logger.InfoContext(ctx, "read item", "ids", ids)

		return response, nil
	}
}
