package utils

import "errors"

var ErrMissingCredentials = errors.New("missing credentials")
var ErrSomethingWentWrong = errors.New("something went wrong, try again later")
