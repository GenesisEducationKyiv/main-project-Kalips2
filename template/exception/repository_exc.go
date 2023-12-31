package exception

import "github.com/pkg/errors"

var (
	ErrReadFromStorage                 = errors.New("Failed to read from storage")
	ErrWriteToStorage                  = errors.New("Failed to write to storage")
	ErrJsonWithIncorrectFormat         = errors.New("Failed to parse file due to incorrect format")
	FailToAddEmailToStorageMessage     = "Failed to subscribe email"
	FailToCheckExistenceOfEmailMessage = "Failed to check the existence of email"
	FailToGetEmailsMessage             = "Failed to get emails from storage"
)
