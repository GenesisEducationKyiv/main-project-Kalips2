package provider

import "time"

type DefaultTimeProvider struct{}

func NewDefaultTimeProvider() *DefaultTimeProvider {
	return &DefaultTimeProvider{}
}

func (t *DefaultTimeProvider) Now() time.Time {
	return time.Now()
}
