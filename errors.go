package sourcify

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"io"
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

func ToErrorResponse(response io.ReadCloser) error {
	var errorResp ErrorResponse
	if err := json.NewDecoder(response).Decode(&errorResp); err == nil && errorResp.Message != "" {
		return fmt.Errorf("sourcify returned error (%s): %s", errorResp.CustomCode, errorResp.Message)
	}
	return nil
}
