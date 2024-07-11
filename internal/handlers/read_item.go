package handlers

import (
	"context"
	"encoding/json"
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
type readItemRequest struct {
	ID string `json:"id"`
}

func (req readItemRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)
	if req.ID == "" {
		problems["id"] = "missing required property id"
	}
	return problems
}

// readItemResponse represents a response for reading an item. The lambda
// function is confiured for batch invoke so the response is a batch invoke
// response.
type readItemResponse = batchInvokeResponse[models.Item]

// HandleReadItem is the handler for the ReadItem Lambda function.
func HandleReadItem(
	config HandleReadItemConfig,
	logger *slog.Logger,
	dynamodbclient handleReadItemDynamoDBClient,
) func(ctx context.Context, reqs []request[readItemRequest]) ([]readItemResponse, error) {
	return func(ctx context.Context, reqs []request[readItemRequest]) ([]readItemResponse, error) {
		responses := make([]readItemResponse, len(reqs))

		logger.InfoContext(ctx, "read item", "request", reqs)

		for i, req := range reqs {
			if problems := req.Valid(ctx); len(problems) > 0 {
				responses[i].Error = errInvalidRequest[readItemRequest]{ctx, req}
			}
		}

		b, _ := json.MarshalIndent(responses, "", "\t")

		logger.InfoContext(ctx, "responses", "responses", responses, "raw", string(b))
		return responses, nil
	}
}
