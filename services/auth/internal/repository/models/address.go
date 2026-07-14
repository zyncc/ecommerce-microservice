package models

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FirstName string
	LastName  *string
	Email     string
	Phone     string
	Address1  string
	Address2  *string
	City      string
	State     string
	Zip       string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateAddressParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FirstName string
	LastName  *string
	Email     string
	Phone     string
	Address1  string
	Address2  *string
	City      string
	State     string
	Zip       string
}
