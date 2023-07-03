package exception

import "github.com/pkg/errors"

var (
	ErrEmailIsAlreadySubscribed = errors.New("Email is already subscribed!")
)
