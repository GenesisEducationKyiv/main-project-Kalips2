package error

import "github.com/pkg/errors"

var (
	ErrReadFromStorage         = errors.New("Failed to read from storage")
	ErrWriteToStorage          = errors.New("Failed to write to storage")
	ErrJsonWithIncorrectFormat = errors.New("Failed to parse file due to incorrect format")
)
