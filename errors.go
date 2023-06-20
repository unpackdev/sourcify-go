package sourcify

import "fmt"

var (
	// ErrInvalidParamType represents an error when a parameter of an invalid type is encountered.
	ErrInvalidParamType = func(t string) error {
		return fmt.Errorf("encountered a parameter of invalid type: %s", t)
	}
)
