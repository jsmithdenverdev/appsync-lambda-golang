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
type readItemResponse struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
}

// HandleReadItem is the handler for the ReadItem Lambda function. The handler
// is batch invoked. This means it recieves an array of requests and needs to
// return an array of results with the same length and positions as the input.
func HandleReadItem(
	config HandleReadItemConfig,
	logger *slog.Logger,
	dynamodbclient handleReadItemDynamoDBClient,
) func(ctx context.Context, reqs []request[readItemRequest]) ([]readItemResponse, []error) {
	return func(ctx context.Context, reqs []request[readItemRequest]) ([]readItemResponse, []error) {
		var (
			responses = make([]readItemResponse, len(reqs))
			keys      = make([]map[string]types.AttributeValue, len(reqs))
			errors    = make([]error, len(reqs))
		)

		logger.InfoContext(ctx, "[query: item]", "request", reqs)

		for i, req := range reqs {
			if problems := req.Valid(ctx); len(problems) > 0 {
				errors[i] = errInvalidRequest[readItemRequest]{ctx, req}
			} else {
				// Assign a key to look this item up
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

		results, err := dynamodbclient.BatchGetItem(ctx, batchGetInput)
		if err != nil {
			logger.ErrorContext(ctx, "failed to batch get items", "error", err)
			for i, _ := range errors {
				errors[i] = errInternalServer
			}
			return responses, errors
		}

		for _, result := range results.Responses[config.TableName] {
			var item models.Item

			if err := attributevalue.UnmarshalMap(result, &item); err != nil {
				// We failed to unmarshal, but we might still be able to extract an ID.
				// This is important because we must return results and errors in the
				// same order as the original requests.
				// DynamoDB doesn't guarantee the ordering of its results, so we can't
				// rely on the index of this record in the result set. Instead we try
				// to find a common value to match this error to a request. The value
				// we'll use is id.
				idAttributeValue := result["id"].(*types.AttributeValueMemberS)
				logger.ErrorContext(ctx, "failed to unmarshal item", "error", err)

				// We found an id in the DynamoDB record.
				if idAttributeValue != nil && idAttributeValue.Value != "" {
					var (
						index int  = 0
						match bool = false
					)
					// Now we need to attempt to find the index of the request with the
					// matching ID.
					for i, request := range reqs {
						// We matched this error with its corresponding request and have an
						// index to work with
						if idAttributeValue.Value == request.Args.ID {
							index = i
							match = true
							break
						}
					}

					if match {
						errors[index] = errInternalServer
						continue
					} else {
						// If we couldn't match this unmarshalling error back to the same
						// index as a request we can no longer rely on ordering and need to
						// return an error for each request.
						for i := range errors {
							errors[i] = errInternalServer
						}
						break
					}

				} else {
					// If we couldn't find an ID, we have no way of tying this error back
					// to a particular request, and need to attach an error to every
					// request index.
					for i := range errors {
						errors[i] = errInternalServer
					}
					break
				}
			}

			// We successfully unmarshalled the record, and now we need to find the
			// index of the corresponding request
			var (
				index int  = 0
				match bool = false
			)
			// Now we need to attempt to find the index of the request with the
			// matching ID.
			for i, request := range reqs {
				// We matched this error with its corresponding request and have an
				// index to work with
				if item.ID == request.Args.ID {
					index = i
					match = true
					break
				}
			}

			if match {
				responses[index] = readItemResponse{
					ID:   item.ID,
					Name: item.Name,
					Tags: item.Tags,
				}
			}
		}

		return responses, errors
	}
}
