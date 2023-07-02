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
	conf            config.CryptoConfig
	rateService     handler.RateService
	emailRepository EmailRepository
	emailSender     GoMailSender
}

type EmailRepository interface {
	SaveEmail(email model.Email) error
	GetEmailsFromStorage() ([]model.Email, error)
	CheckEmailIsExist(email model.Email) (bool, error)
}

type GoMailSender interface {
	SendMessageTo(message *model.Message, recipients []string) error
}

func (emailService *EmailServiceImpl) SendRateToEmails() error {
	var emails []model.Email
	var err error
	conf := emailService.conf

	emails, err = emailService.emailRepository.GetEmailsFromStorage()
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	rate, err := emailService.rateService.GetRate()
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

func (emailService *EmailServiceImpl) SubscribeEmail(emailVal string) error {
	var err error

	err = validateEmail(emailVal)
	if err != nil {
		return errors.Wrap(err, message.FailToSubscribeMessage)
	}

	email := model.Email{Mail: emailVal}
	exist, err := emailService.emailRepository.CheckEmailIsExist(email)
	if exist {
		err = exception.ErrEmailIsAlreadySubscribed
	}
	if err != nil {
		return errors.Wrap(err, message.FailToSubscribeMessage)
	}

	err = emailService.emailRepository.SaveEmail(email)
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

func NewEmailService(c config.CryptoConfig, service handler.RateService, emailRepository EmailRepository, sender GoMailSender) *EmailServiceImpl {
	return &EmailServiceImpl{
		conf:            c,
		rateService:     service,
		emailRepository: emailRepository,
		emailSender:     sender,
	}
}
