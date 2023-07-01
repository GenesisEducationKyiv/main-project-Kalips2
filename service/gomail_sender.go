package service

import (
	"btc-app/config"
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"strings"
)

type GoMailSender interface {
	CreateMessage(emailFrom string, header string, body string) *gomail.Message
	SendMessageTo(message *gomail.Message, recipients []string) error
}

type GoMailSenderImpl struct {
	dialer *gomail.Dialer
}

func (sender *GoMailSenderImpl) CreateMessage(emailFrom string, header string, body string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", emailFrom)
	message.SetHeader("Subject", header)
	message.SetBody("text/plain", body)
	return message
}

func (sender *GoMailSenderImpl) SendMessageTo(message *gomail.Message, recipients []string) error {
	var err error
	var failedEmails []string

	for _, email := range recipients {
		message.SetHeader("To", email)
		if err = sender.dialer.DialAndSend(message); err != nil {
			failedEmails = append(failedEmails, email)
		}
	}
	return errors.Wrap(err, "Failed to send emails to: "+strings.Join(failedEmails, " "))
}

func NewEmailSender(c *config.Config) *GoMailSenderImpl {
	dialer := gomail.NewDialer(c.EmailServiceHost, c.EmailServicePort, c.EmailServiceFrom, c.EmailServicePassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &GoMailSenderImpl{dialer: dialer}
}
