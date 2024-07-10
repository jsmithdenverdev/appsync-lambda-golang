package handlers

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/models"
)

// HandleListItemsConfig represents the config needed by the ListItems Lambda
// function.
type HandleListItemsConfig struct {
	TableName string
}

// handleListItemsDynamoDBClient represents the Dynamo DB operations that can
// be performed from HandleCreateItem.
type handleListItemsDynamoDBClient interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// listItemsRequest represents a request to read an item. The lambda function is
// confiured for batch invoke so the request is a slice of ids instead of a
// single id.
type listItemsRequest []string

func (req listItemsRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)
	return problems
}

// listItemsRequest represents a request to read an item. The lambda function is
// confiured for batch invoke so the request is a slice of ids instead of a
// single id.
type listItemsResponse []models.Item

// HandleCreateItem is the handler for the CreateItem Lambda function.
// The handler consumes a `requests.CreateItem` request and attempts to create
// a new `models.Item` record, before saving that to the items DynamoDB table.
// After the record is peristed, the full record is read from the table before
// being marshalled into a `responses.Item` response.
// If an error occurs contextual logging is done around the error and it is
// bubbled up the stack.
func HandleListItems(
	config HandleListItemsConfig,
	logger *slog.Logger,
	dynamodbclient handleListItemsDynamoDBClient,
) func(ctx context.Context, req request[listItemsRequest]) (listItemsResponse, error) {
	return func(ctx context.Context, req request[listItemsRequest]) (listItemsResponse, error) {
		var response listItemsResponse

		return response, nil
	}
}
