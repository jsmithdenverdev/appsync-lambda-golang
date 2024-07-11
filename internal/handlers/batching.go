package handlers

import "encoding/json"

type batchInvokeResponse[T any] struct {
	Value T     `json:"value"`
	Error error `json:"error"`
}

// MarshalJSON implements the json.Marshaler interface for batchInvokeResponse
func (b batchInvokeResponse[T]) MarshalJSON() ([]byte, error) {
	// type Alias batchInvokeResponse[T] // Create an alias to avoid infinite recursion

	// Helper function to convert error to string
	errorToString := func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	}

	return json.Marshal(&struct {
		Value T      `json:"value"`
		Error string `json:"error"`
	}{
		Value: b.Value,
		Error: errorToString(b.Error),
	})
}
