package application

import (
	"btc-app/config"
	"btc-app/pkg/domain/model"
	"btc-app/pkg/domain/service"
	cerror "btc-app/template/cerror"
	"btc-app/template/message"
	"github.com/pkg/errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type (
	EmailServiceImpl struct {
		conf            config.CryptoConfig
		rateService     service.RateService
		emailRepository EmailRepository
		emailSender     GoMailSender
		logger          service.Logger
	}

	EmailRepository interface {
		SaveEmail(email model.Email) error
		GetEmailsFromStorage() ([]model.Email, error)
		CheckEmailIsExist(email model.Email) (bool, error)
	}

	GoMailSender interface {
		SendMessageTo(message *model.EmailMessage, recipients []model.Email) error
	}
)

func (emailService *EmailServiceImpl) SendRateToEmails() error {
	var emails []model.Email
	var err error
	conf := emailService.conf

	emails, err = emailService.emailRepository.GetEmailsFromStorage()
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	curPair := model.NewCurrencyPair(emailService.conf.CurrencyTo, emailService.conf.CurrencyFrom)
	rate, err := emailService.rateService.GetCurrencyRate(*curPair)
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	msg := model.NewRateMessage(rate, conf.CurrencyFrom, conf.CurrencyTo)
	err = emailService.emailSender.SendMessageTo(msg, emails)
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}
	return err
}

func (emailService *EmailServiceImpl) SubscribeEmail(email model.Email) error {
	var err error

	err = validateEmail(email)
	if err != nil {
		return errors.Wrap(err, message.FailToSubscribeMessage)
	}

	exist, err := emailService.emailRepository.CheckEmailIsExist(email)
	if exist {
		err = cerror.ErrEmailIsAlreadySubscribed
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

func validateEmail(email model.Email) error {
	var err error
	if !emailRegex.MatchString(email.GetAddress()) {
		err = errors.New("Email doesn't match regex: " + emailRegex.String())
	}
	return err
}

func NewEmailService(c config.CryptoConfig, rateService service.RateService, emailRepository EmailRepository, sender GoMailSender) *EmailServiceImpl {
	return &EmailServiceImpl{
		conf:            c,
		rateService:     rateService,
		emailRepository: emailRepository,
		emailSender:     sender,
	}
}
