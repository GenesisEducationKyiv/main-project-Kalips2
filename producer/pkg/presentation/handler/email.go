package handler

import (
	"net/http"
	"producer/config"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
	"producer/pkg/presentation/presenter"
)

type EmailHandlerImpl struct {
	conf         *config.Config
	emailService service.EmailService
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
		email := model.NewEmail(r.FormValue("email"))
		if err := emailHr.emailService.SubscribeEmail(email); err != nil {
			presenter.PresentErrorByConflict(w, err)
		} else {
			presenter.PresentString(w, "Emails have been subscribed.")
		}
	}
}

func NewEmailHandler(c *config.Config, emailService service.EmailService) *EmailHandlerImpl {
	return &EmailHandlerImpl{
		conf:         c,
		emailService: emailService,
	}
}
