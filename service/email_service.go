package service

import (
	"btc-app/config"
	"btc-app/repository"
	"crypto/tls"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var (
	emailRegex                  = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	ErrEmailIsAlreadySubscribed = errors.New("Email is already subscribed!")
	failToSendRateMessage       = "Failed to send the rate to emails"
	failToSubscribeMessage      = "Failed to subscribe email"
)

type EmailServiceImpl struct {
	conf            *config.Config
	emailRepository EmailRepository
}

type EmailRepository interface {
	SaveEmailToStorage(email string) error
	GetEmailsFromStorage() ([]string, error)
}

func (emailService EmailServiceImpl) SendRateToEmails() error {
	var emails []string
	var err error

	emails, err = emailService.emailRepository.GetEmailsFromStorage()
	if err != nil {
		return errors.Wrap(err, failToSendRateMessage)
	}

	rate, err := GetCurrentRate(c)
	if err != nil {
		return errors.Wrap(err, failToSendRateMessage)
	}

	dialer, message := setUpMessageToSend(rate, emailService.conf)
	err = sendMessageToEmails(message, emails, dialer)
	if err != nil {
		return errors.Wrap(err, failToSendRateMessage)
	}
	return err
}

func (emailService EmailServiceImpl) SubscribeEmail(email string) error {
	var err error

	err = validateEmail(email)
	if err != nil {
		return errors.Wrap(err, failToSubscribeMessage)
	}

	exist, err := repository.CheckEmailIsExist(email, c.EmailStoragePath)
	if exist {
		err = ErrEmailIsAlreadySubscribed
	}
	if err != nil {
		return errors.Wrap(err, failToSubscribeMessage)
	}

	err = repository.SaveEmailToStorage(email, c.EmailStoragePath)
	if err != nil {
		return errors.Wrap(err, failToSubscribeMessage)
	}
	return err
}

func validateEmail(email string) error {
	var err error
	if !emailRegex.MatchString(email) {
		err = errors.New("Email doesn't match regex: " + emailRegex.String())
	}
	return err
}

func setUpMessageToSend(rate float64, c *config.Config) (*gomail.Dialer, *gomail.Message) {
	dialer := gomail.NewDialer(c.EmailServiceHost, c.EmailServicePort, c.EmailServiceFrom, c.EmailServicePassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	message := gomail.NewMessage()
	message.SetHeader("From", c.EmailServiceFrom)
	message.SetHeader("Subject", c.EmailServiceSubject)
	message.SetBody("text/plain", "Поточний курс "+c.CurrencyFrom+" до "+c.CurrencyTo+": "+fmt.Sprintf("%.5f", rate)+".")

	return dialer, message
}

func sendMessageToEmails(message *gomail.Message, emails []string, dialer *gomail.Dialer) error {
	var err error
	var failedEmails []string
	for _, email := range emails {
		message.SetHeader("To", email)
		if err = dialer.DialAndSend(message); err != nil {
			failedEmails = append(failedEmails, email)
		}
	}
	return errors.Wrap(err, "Failed to send emails to: "+strings.Join(failedEmails, " "))
}
func NewEmailService(c *config.Config) EmailServiceImpl {
	return EmailServiceImpl{
		conf:            c,
		emailRepository: repository.NewEmailRepository(c.EmailStoragePath),
	}
}
