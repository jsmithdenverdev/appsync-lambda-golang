package handlers

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/models"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/requests"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/responses"
)

// HandleCreateItemConfig represents the config needed by the CreateItem Lambda
// function.
type HandleCreateItemConfig struct {
	TableName string
}

// handleCreateItemDynamoDBClient represents the Dynamo DB operations that can
// be performed from HandleCreateItem.
type handleCreateItemDynamoDBClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// HandleCreateItem is the handler for the CreateItem Lambda function.
// The handler consumes a `requests.CreateItem` request and attempts to create
// a new `models.Item` record, before saving that to the items DynamoDB table.
// After the record is peristed, the full record is read from the table before
// being marshalled into a `responses.Item` response.
// If an error occurs contextual logging is done around the error and it is
// bubbled up the stack.
func HandleCreateItem(
	config HandleCreateItemConfig,
	logger *slog.Logger,
	dynamodbclient handleCreateItemDynamoDBClient,
) func(ctx context.Context, request requests.CreateItem) (responses.CreateItem, error) {
	return func(ctx context.Context, request requests.CreateItem) (responses.CreateItem, error) {
		var response responses.CreateItem

		if problems := request.Valid(ctx); len(problems) > 0 {
			return response, validationError{problems}
		}

		id := uuid.New().String()

		av, err := attributevalue.MarshalMap(models.Item{
			ID:   id,
			Name: request.Name,
		})

		if err != nil {
			logger.ErrorContext(ctx, "failed to marshal item to attribute value map", "error", err)
			return response, errInternalServer
		}

		_, err = dynamodbclient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(config.TableName),
			Item:      av,
		})

		if err != nil {
			logger.ErrorContext(ctx, "failed to put item", "error", err)
			return response, errInternalServer
		}

		row, err := dynamodbclient.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(config.TableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{
					Value: id,
				},
			},
		})

		if err != nil {
			logger.ErrorContext(ctx, "failed to get item", "error", err)
			return response, errInternalServer
		}

		var item responses.Item

		if err = attributevalue.UnmarshalMap(row.Item, &item); err != nil {
			logger.ErrorContext(ctx, "failed to marshal item from attribute value map", "error", err)
			return response, errInternalServer
		}

		return responses.CreateItem{
			Item: item,
		}, nil
	}
}
