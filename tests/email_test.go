package tests

import (
	"btc-app/config"
	"testing"
)

func TestEmailWasSubscribed(t *testing.T) {
	var c config.Config
	c.EmailStoragePath("subscriptions-test.csv")
}
