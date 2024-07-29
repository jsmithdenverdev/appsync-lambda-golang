package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServer = errors.New(http.StatusText(http.StatusInternalServerError))
)
