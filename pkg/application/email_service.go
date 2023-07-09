package application

import (
	"btc-app/config"
	"btc-app/pkg/domain"
	cerror "btc-app/template/cerror"
	"btc-app/template/message"
	"github.com/pkg/errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type (
	EmailServiceImpl struct {
		conf            config.CryptoConfig
		rateService     RateService
		emailRepository EmailRepository
		emailSender     GoMailSender
	}

	RateService interface {
		GetRate(curPair domain.CurrencyPair) (*domain.CurrencyRate, error)
	}

	EmailService interface {
		SendRateToEmails() error
		SubscribeEmail(emailVal string) error
	}

	EmailRepository interface {
		SaveEmail(email domain.Email) error
		GetEmailsFromStorage() ([]domain.Email, error)
		CheckEmailIsExist(email domain.Email) (bool, error)
	}

	GoMailSender interface {
		SendMessageTo(message *domain.EmailMessage, recipients []domain.Email) error
	}
)

func (emailService *EmailServiceImpl) SendRateToEmails() error {
	var emails []domain.Email
	var err error
	conf := emailService.conf

	emails, err = emailService.emailRepository.GetEmailsFromStorage()
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	curPair := domain.NewCurrencyPair(emailService.conf.CurrencyTo, emailService.conf.CurrencyFrom)
	rate, err := emailService.rateService.GetRate(*curPair)
	if err != nil {
		return errors.Wrap(err, message.FailToSendRateMessage)
	}

	msg := domain.NewRateMessage(rate, conf.CurrencyFrom, conf.CurrencyTo)
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

	email := domain.NewEmail(emailVal)
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

func validateEmail(email string) error {
	var err error
	if !emailRegex.MatchString(email) {
		err = errors.New("Email doesn't match regex: " + emailRegex.String())
	}
	return err
}

func NewEmailService(c config.CryptoConfig, rateService RateService, emailRepository EmailRepository, sender GoMailSender) *EmailServiceImpl {
	return &EmailServiceImpl{
		conf:            c,
		rateService:     rateService,
		emailRepository: emailRepository,
		emailSender:     sender,
	}
}
