package server

import (
	"btc-app/config"
	"btc-app/handler"
	"btc-app/repository"
	"btc-app/service"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type Server struct {
	Config *config.Config
	Router *chi.Mux
}

type EmailHandler interface {
	SendToEmailsHandler() http.HandlerFunc
	SubscribeEmailHandler() http.HandlerFunc
}

type RateHandler interface {
	GetCurrentRateHandler() http.HandlerFunc
}

func (s *Server) SetupServer(c *config.Config) {
	emailRepository := repository.NewEmailRepository(c.Database)
	emailSender := service.NewEmailSender(c.MailService)

	rateService := service.NewRateService(c.Crypto, service.NewChainOfProviders(c.Crypto))
	emailService := service.NewEmailService(c.Crypto, rateService, emailRepository, emailSender)

	rateHandler := handler.NewRateHandler(c, rateService)
	emailHandler := handler.NewEmailHandler(c, emailService)

	s.Router.Get("/rate", rateHandler.GetCurrentRateHandler())
	s.Router.Post("/subscribe", emailHandler.SubscribeEmailHandler())
	s.Router.Post("/sendEmails", emailHandler.SendToEmailsHandler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port), s.Router); err != nil {
		log.Fatalf(fmt.Sprintf("Failed to listen and serve server on the port - %d", c.Server.Port), err)
	}
}

func NewServer(conf *config.Config) *Server {
	return &Server{
		Config: conf,
		Router: chi.NewRouter(),
	}
}
