package model

import (
	"fmt"
	"time"
)

var (
	ErrInvalidReqData = NewAppError("invalid_request_data", "Invalid request data", ErrorTypeRequest)

	ErrorTypeValidation = "validation"
	ErrorTypeDatabase   = "database"
	ErrorTypeRequest    = "request"
	ErrorTypeServer     = "server"
)

type AppError struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Type      string    `json:"-"`
	Timestamp time.Time `json:"timestamp"`
}

func NewAppError(code string, message string, errorType string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Type:      errorType,
		Timestamp: time.Now(),
	}
}

func (err *AppError) Error() string {
	return fmt.Sprintf("(%s) [%s]: %s", err.Type, err.Code, err.Message)
}
