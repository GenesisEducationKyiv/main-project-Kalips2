package service

import (
	"btc-app/config"
	"btc-app/handler"
	"btc-app/model"
	"btc-app/template/exception"
	"btc-app/template/message"
	"github.com/pkg/errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type EmailServiceImpl struct {
	conf            *config.Config
	rateService     handler.RateService
	emailRepository EmailRepository
	emailSender     GoMailSender
}

type EmailRepository interface {
	SaveEmailToStorage(email string) error
	GetEmailsFromStorage() ([]string, error)
	CheckEmailIsExist(email string) (bool, error)
}

type GoMailSender interface {
	SendMessageTo(message *model.Message, recipients []string) error
}

func (emailService *EmailServiceImpl) SendRateToEmails() error {
	var emails []string
	var err error
	conf := emailService.conf

	emails, err = emailService.emailRepository.GetEmailsFromStorage()
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	rate, err := emailService.rateService.GetCurrentRate()
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	msg := model.NewRateMessage(rate, conf.EmailServiceFrom, conf.CurrencyFrom, conf.CurrencyTo)
	err = emailService.emailSender.SendMessageTo(msg, emails)
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}
	return err
}

func (emailService *EmailServiceImpl) SubscribeEmail(email string) error {
	var err error

	err = validateEmail(email)
	if err != nil {
		return errors.Wrap(err, message.FailToSubscribeMessage)
	}

	exist, err := emailService.emailRepository.CheckEmailIsExist(email)
	if exist {
		err = exception.ErrEmailIsAlreadySubscribed
	}
	if err != nil {
		return errors.Wrap(err, message.FailToSubscribeMessage)
	}

	err = emailService.emailRepository.SaveEmailToStorage(email)
	if err != nil {
		return errors.Wrap(err, message.FailToSubscribeMessage)
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

func NewEmailService(c *config.Config, service handler.RateService, emailRepository EmailRepository, sender GoMailSender) *EmailServiceImpl {
	return &EmailServiceImpl{
		conf:            c,
		rateService:     service,
		emailRepository: emailRepository,
		emailSender:     sender,
	}
}
