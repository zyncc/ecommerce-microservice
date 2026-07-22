package types

import "errors"

var (
	ErrInventoryNotFound = errors.New("inventory not found")
	ErrInsufficientStock = errors.New("quantity is out of stock")
	ErrDatabase          = errors.New("database error")
	ErrInvalidSize       = errors.New("size does not exist")
)
