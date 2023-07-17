package service

import "producer/pkg/domain/model"

type EmailService interface {
	SendRateToEmails() error
	SubscribeEmail(email model.Email) error
}
