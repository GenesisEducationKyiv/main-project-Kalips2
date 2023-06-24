package service

import (
	"btc-app/config"
	"btc-app/repository"
	"btc-app/template/exception"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type EmailService interface {
	SendRateToEmails() error
	SubscribeEmail(email string) error
}

type EmailServiceImpl struct {
	conf            *config.Config
	rateService     RateService
	emailRepository repository.EmailRepository
	emailSender     GoMailSender
}

func (emailService *EmailServiceImpl) SendRateToEmails() error {
	var emails []string
	var err error
	conf := emailService.conf

	emails, err = emailService.emailRepository.GetEmailsFromStorage()
	if err != nil {
		return errors.Wrap(err, exception.FailToSendRateMessage)
	}

	rate, err := emailService.rateService.GetCurrentRate()
	if err != nil {
		return errors.Wrap(err, exception.FailToSendRateMessage)
	}

	rateFormatted := strconv.FormatFloat(rate, 'f', 5, 64)
	emailSubject := "Поточний курс " + conf.CurrencyFrom + " до " + conf.CurrencyTo + "."
	message := emailService.emailSender.CreateMessage(conf.EmailServiceFrom, emailSubject, rateFormatted)

	err = emailService.emailSender.SendMessageTo(message, emails)
	if err != nil {
		return errors.Wrap(err, exception.FailToSendRateMessage)
	}
	return err
}

func (emailService *EmailServiceImpl) SubscribeEmail(email string) error {
	var err error

	err = validateEmail(email)
	if err != nil {
		return errors.Wrap(err, exception.FailToSubscribeMessage)
	}

	exist, err := emailService.emailRepository.CheckEmailIsExist(email)
	if exist {
		err = exception.EmailIsAlreadySubscribed
	}
	if err != nil {
		return errors.Wrap(err, exception.FailToSubscribeMessage)
	}

	err = emailService.emailRepository.SaveEmailToStorage(email)
	if err != nil {
		return errors.Wrap(err, exception.FailToSubscribeMessage)
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

func (emailService *EmailServiceImpl) SetEmailRepository(repo repository.EmailRepository) {
	emailService.emailRepository = repo
}

func (emailService *EmailServiceImpl) SetEmailSender(sender GoMailSender) {
	emailService.emailSender = sender
}

func (emailService *EmailServiceImpl) SetRateService(rateService RateService) {
	emailService.rateService = rateService
}

func NewEmailService(c *config.Config) *EmailServiceImpl {
	return &EmailServiceImpl{
		conf:            c,
		rateService:     NewRateService(c),
		emailRepository: repository.NewEmailRepository(c.EmailStoragePath),
		emailSender:     NewEmailSender(c),
	}
}
