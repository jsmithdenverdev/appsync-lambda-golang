package handlers

import (
	"context"
)

// A request represents a GraphQL request that has been mapped from an AWS
// AppSync API into a standardized input template for Lambda functions using
// Velocity Template Language
//
// The VTL for a request looks like so. The full arguments from the GraphQL
// request are attached to an args property in the payload.
//
// ```
//
//	{
//	  "version": "2018-05-29",
//	  "operation": "Invoke" | "BatchInvoke",
//	  "payload": {
//	      "args": $util.toJson($ctx.args)
//	  }
//	}
//
// ```
// request implements validator and calls the validate method of its args
// property allowing enforcing request validation to be implemented for all
// requests.
type request[T validator] struct {
	// Args from the original GraphQL request.
	//
	// Args implements validator and can be validated.
	Args T `json:"args"`
}

// Valid validates the args of the request.
func (req request[T]) Valid(ctx context.Context) (problems map[string]string) {
	return req.Args.Valid(ctx)
}
