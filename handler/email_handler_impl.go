package handler

import (
	"btc-app/config"
	"btc-app/service"
	"fmt"
	"net/http"
)

type EmailHandlerImpl struct {
	conf         *config.Config
	emailService EmailService
}

type EmailService interface {
	SendRateToEmails() error
	SubscribeEmail(email string) error
}

func (emailHr EmailHandlerImpl) SendToEmailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := emailHr.emailService.SendRateToEmails(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Emails have been sent.")
		}
	}
}

func (emailHr EmailHandlerImpl) SubscribeEmailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")

		if err := emailHr.emailService.SubscribeEmail(email); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Emails have been subscribed.")
		}
	}
}

func NewEmailHandler(c *config.Config) *EmailHandlerImpl {
	return &EmailHandlerImpl{
		emailService: service.NewEmailService(c),
		conf:         c,
	}
}
