package resolvers

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	// A generic internal server error.
	errInternalServer = errors.New("an error occured")
)

// An error stemming from a request that returned problems from Valid.
type errInvalidRequest[T validator] struct {
	ctx     context.Context
	request request[T]
}

// Error returns the problems for the request.
func (err errInvalidRequest[T]) Error() string {
	errstr := "the following problems were detected: %s"
	errstrs := make([]string, 0)
	for k, v := range err.request.Valid(err.ctx) {
		errstrs = append(errstrs, fmt.Sprintf("%s: %s", k, v))
	}
	return fmt.Sprintf(errstr, strings.Join(errstrs, ", "))
}
