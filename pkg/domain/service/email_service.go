package service

import "btc-app/pkg/domain/model"

type EmailService interface {
	SendRateToEmails() error
	SubscribeEmail(email model.Email) error
}
