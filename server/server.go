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

func NewServer(conf config.Config) *Server {
	return &Server{
		Config: conf,
		Router: chi.NewRouter(),
	}
}

func (s *Server) InitHandlers(c *config.Config) {
	s.Router.Get("/rate", func(w http.ResponseWriter, r *http.Request) {
		handler.GetCurrentRateHandler(w, r, &s.Config)
	})

	s.Router.Post("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		handler.SubscribeEmailHandler(w, r, &s.Config)
	})

	s.Router.Post("/sendEmails", func(w http.ResponseWriter, r *http.Request) {
		handler.SendToEmailsHandler(w, r, &s.Config)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), s.Router); err != nil {
		log.Fatalf(fmt.Sprintf("Failed to listen and serve server on the port - %d", c.Port), err)
	}
}
