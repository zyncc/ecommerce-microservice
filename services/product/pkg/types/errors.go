package types

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrDatabase = errors.New("database error")
var ErrSendKafkaMessage = errors.New("failed to send kafka message")
