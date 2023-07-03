package service

import (
	"btc-app/config"
	"btc-app/model"
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"strings"
)

type GoMailSenderImpl struct {
	dialer *gomail.Dialer
}

func (sender *GoMailSenderImpl) createMailMessage(message *model.Message) *gomail.Message {
	mailMsg := gomail.NewMessage()
	mailMsg.SetHeader("From", message.EmailFrom)
	mailMsg.SetHeader("Subject", message.Header)
	mailMsg.SetBody("text/plain", message.Body)
	return mailMsg
}

func (sender *GoMailSenderImpl) SendMessageTo(message *model.Message, recipients []string) error {
	var err error
	var failedEmails []string

	mailMsg := sender.createMailMessage(message)
	for _, email := range recipients {
		mailMsg.SetHeader("To", email)
		if err = sender.dialer.DialAndSend(mailMsg); err != nil {
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
