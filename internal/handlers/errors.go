package handlers

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInternalServer = errors.New("an error occured - please try again")
)

type validationError struct {
	problems map[string]string
}

func (err validationError) Error() string {
	errstr := "the following problems were detected: %s"
	errstrs := make([]string, 0)
	for k, v := range err.problems {
		errstrs = append(errstrs, fmt.Sprintf("%s: %s", k, v))
	}
	return fmt.Sprintf(errstr, strings.Join(errstrs, ", "))
}
