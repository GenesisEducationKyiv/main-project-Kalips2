package server

import (
	"btc-app/config"
	"btc-app/handler"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type Server struct {
	Config config.Config
	Router *chi.Mux
}

type EmailHandler interface {
	SendToEmailsHandler() http.HandlerFunc
	SubscribeEmailHandler() http.HandlerFunc
}

type RateHandler interface {
	GetCurrentRateHandler() http.HandlerFunc
}

func (s *Server) InitHandlers(c *config.Config) {
	rateHandler := handler.NewRateHandler(c)
	emailHandler := handler.NewEmailHandler(c)

	s.Router.Get("/rate", rateHandler.GetCurrentRateHandler())
	s.Router.Post("/subscribe", emailHandler.SubscribeEmailHandler())
	s.Router.Post("/sendEmails", emailHandler.SendToEmailsHandler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), s.Router); err != nil {
		log.Fatalf(fmt.Sprintf("Failed to listen and serve server on the port - %d", c.Port), err)
	}
}

func NewServer(conf config.Config) *Server {
	return &Server{
		Config: conf,
		Router: chi.NewRouter(),
	}
}
