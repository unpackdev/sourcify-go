package sourcify

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	// ErrInvalidParamType represents an error when a parameter of an invalid type is encountered.
	ErrInvalidParamType = func(t string) error {
		return fmt.Errorf("encountered a parameter of invalid type: %s", t)
	}
)

type ErrorResponse struct {
	ErrorId    uuid.UUID `json:"errorId"`
	CustomCode string    `json:"customCode"`
	Message    string    `json:"message"`
}
