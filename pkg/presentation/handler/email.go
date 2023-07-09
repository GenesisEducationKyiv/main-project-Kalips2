package handler

import (
	"btc-app/config"
	"btc-app/pkg/application"
	"btc-app/pkg/presentation/presenter"
	"net/http"
)

type EmailHandlerImpl struct {
	conf         *config.Config
	emailService application.EmailService
}

func (emailHr *EmailHandlerImpl) SendToEmailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := emailHr.emailService.SendRateToEmails(); err != nil {
			presenter.PresentErrorByInternalServerError(w, err)
		} else {
			presenter.PresentString(w, "Emails have been sent.")
		}
	}
}

func (emailHr *EmailHandlerImpl) SubscribeEmailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")

		if err := emailHr.emailService.SubscribeEmail(email); err != nil {
			presenter.PresentErrorByConflict(w, err)
		} else {
			presenter.PresentString(w, "Emails have been subscribed.")
		}
	}
}

func NewEmailHandler(c *config.Config, emailService application.EmailService) *EmailHandlerImpl {
	return &EmailHandlerImpl{
		conf:         c,
		emailService: emailService,
	}
}
