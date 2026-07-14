package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type Success[T any] struct {
	Code      int       `json:"code"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"ts"`
}

type Error struct {
	Code      int       `json:"code"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"ts"`
}

type ValidationError struct {
	Code      int                 `json:"code"`
	Success   bool                `json:"success"`
	Message   string              `json:"message"`
	Errors    map[string][]string `json:"errors"`
	Timestamp time.Time           `json:"timestamp"`
}

func SuccessResponse[T any](w http.ResponseWriter, code int, message string, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Success[T]{
		Code:      code,
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func ErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Error{
		Code:      code,
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func ValidationErrorResponse(w http.ResponseWriter, errs map[string][]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	json.NewEncoder(w).Encode(ValidationError{
		Code:      http.StatusUnprocessableEntity,
		Success:   false,
		Message:   "validation failed",
		Errors:    errs,
		Timestamp: time.Now(),
	})
}

func AuthorizationErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(Error{
		Code:      http.StatusUnauthorized,
		Success:   false,
		Message:   "Unauthorized: missing or invalid credentials",
		Timestamp: time.Now(),
	})
}

func ForbiddenErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(Error{
		Code:      http.StatusForbidden,
		Success:   false,
		Message:   "Forbidden: you do not have permission to access this resource",
		Timestamp: time.Now(),
	})
}
