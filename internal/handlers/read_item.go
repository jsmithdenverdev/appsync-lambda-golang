package handlers

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/models"
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

// readItemRequest represents a request to read an item. The lambda function is
// confiured for batch invoke so the request is a slice of ids instead of a
// single id.
type readItemRequest []string

// readItemRequest represents a request to read an item. The lambda function is
// confiured for batch invoke so the request is a slice of ids instead of a
// single id.
type readItemResponse []models.Item

func (req readItemRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)
	return problems
}

// HandleReadItem is the handler for the ReadItem Lambda function.
func HandleReadItem(
	config HandleReadItemConfig,
	logger *slog.Logger,
	dynamodbclient handleReadItemDynamoDBClient,
) func(ctx context.Context, req request[readItemRequest]) (readItemResponse, error) {
	return func(ctx context.Context, req request[readItemRequest]) (readItemResponse, error) {
		var response readItemResponse

		if problems := req.Valid(ctx); len(problems) > 0 {
			return response, errInvalidRequest[readItemRequest]{ctx, req}
		}

		logger.InfoContext(ctx, "read item", "request", req)

		return response, nil
	}
}
