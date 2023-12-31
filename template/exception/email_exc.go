package exception

import "github.com/pkg/errors"

var (
	ErrEmailIsAlreadySubscribed = errors.New("Email is already subscribed!")
	FailToSendRateMessage       = "Failed to send the rate to emails"
	FailToSubscribeMessage      = "Failed to subscribe email"
)
