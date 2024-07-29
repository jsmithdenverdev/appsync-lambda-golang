package handlers

import (
	"context"

	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/models"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/services"
)

// createItemRequests represents the fields needed to create a new item.
type CreateItemRequest struct {
	Input struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	} `json:"input"`
}

// Valid checks createItemRequest to ensure it's in a valid state to be used.
func (request CreateItemRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)
	if request.Input.Name == "" {
		problems["input.name"] = "name cannot be empty"
	}
	return problems
}

// CreateItemResponse is the response for a successfully created item.
type CreateItemResponse struct {
	Item models.Item `json:"item"`
}

// HandleCreateItem is the handler for the CreateItem Lambda function.
// The handler consumes a `requests.CreateItem` request and attempts to create
// a new `models.Item` record, before saving that to the items DynamoDB table.
// After the record is peristed, the full record is read from the table before
// being marshalled into a `responses.Item` response.
// If an error occurs contextual logging is done around the error and it is
// bubbled up the stack.
func HandleCreateItem(service services.Item) func(ctx context.Context, req Request[CreateItemRequest]) (CreateItemResponse, error) {
	return func(ctx context.Context, req Request[CreateItemRequest]) (CreateItemResponse, error) {
		var response CreateItemResponse

		if problems := req.Valid(ctx); len(problems) > 0 {
			return response, errInvalidRequest[CreateItemRequest]{ctx, req}
		}

		item, err := service.CreateItem(ctx, services.CreateItemRequest{
			Name: req.Args.Input.Name,
			Tags: req.Args.Input.Tags,
		})

		if err != nil {
			return response, err
		}

		response.Item = item

		return response, nil
	}
}
