package resolvers

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
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
		keys := make([]map[string]types.AttributeValue, len(reqs))

		logger.InfoContext(ctx, "read item", "request", reqs)

		for i, req := range reqs {
			if problems := req.Valid(ctx); len(problems) > 0 {
				responses[i].Error = errInvalidRequest[readItemRequest]{ctx, req}
			} else {
				keys[i] = map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: req.Args.ID},
				}
			}
		}

		batchGetInput := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				config.TableName: {
					Keys: keys,
				},
			},
		}

		result, err := dynamodbclient.BatchGetItem(ctx, batchGetInput)
		if err != nil {
			logger.ErrorContext(ctx, "failed to batch get items", "error", err)
			return responses, errInternalServer
		}

		if len(result.Responses[config.TableName]) != len(reqs) {
			logger.ErrorContext(ctx, "data loader mismatch", "error", err)
			return responses, errInternalServer
		}

		for i, row := range result.Responses[config.TableName] {
			var item models.Item

			if err := attributevalue.UnmarshalMap(row, &item); err != nil {
				logger.ErrorContext(ctx, "failed to unmarshal item", "error", err)
				responses[i] = readItemResponse{Error: errInternalServer}
				continue
			}

			responses[i] = readItemResponse{Value: item}
		}

		return responses, nil
	}
}
