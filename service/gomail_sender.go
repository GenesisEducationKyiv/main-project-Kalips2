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

func (sender *GoMailSenderImpl) createMailMessage(message *model.Message, mailFrom string) *gomail.Message {
	mailMsg := gomail.NewMessage()
	mailMsg.SetHeader("From", mailFrom)
	mailMsg.SetHeader("Subject", message.Header)
	mailMsg.SetBody("text/plain", message.Body)
	return mailMsg
}

func (sender *GoMailSenderImpl) SendMessageTo(message *model.Message, recipients []model.Email) error {
	var err error
	var failedEmails []string

	mailMsg := sender.createMailMessage(message, sender.dialer.Username)
	for _, email := range recipients {
		mailMsg.SetHeader("To", email.Mail)
		if err = sender.dialer.DialAndSend(mailMsg); err != nil {
			failedEmails = append(failedEmails, email.Mail)
		}
	}
	return errors.Wrap(err, "Failed to send emails to: "+strings.Join(failedEmails, ", "))
}

func NewEmailSender(c config.MailConfig) *GoMailSenderImpl {
	dialer := gomail.NewDialer(c.MailHost, c.MailPort, c.MailFrom, c.MailPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &GoMailSenderImpl{dialer: dialer}
}
