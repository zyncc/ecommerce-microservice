package types

import "errors"

var ErrInventoryNotFound = errors.New("inventory not found")
var ErrDatabase = errors.New("database error")
var ErrInvalidSize = errors.New("size does not exist")
