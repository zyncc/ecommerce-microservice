package utils

import "errors"

var ErrMissingCredentials = errors.New("missing credentials")
var ErrSomethingWentWrong = errors.New("something went wrong, try again later")

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}
