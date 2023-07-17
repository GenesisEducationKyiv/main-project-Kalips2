package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"producer/config"
	"producer/pkg/application"
	"producer/pkg/domain/service"
	clogger "producer/pkg/infrastructure/logger"
	"producer/pkg/infrastructure/provider"
	"producer/pkg/infrastructure/repository"
	"producer/pkg/infrastructure/sender"
	"producer/pkg/presentation/handler"
	"time"
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
	logger := setupLogger()
	emailRepository := repository.NewEmailRepository(c.Database, logger)
	emailSender := sender.NewEmailSender(c.MailService)

	rateProvider := setupRateProvider(c, logger)
	rateService := application.NewRateService(c.Crypto, rateProvider, logger)
	emailService := application.NewEmailService(c.Crypto, rateService, emailRepository, emailSender, logger)

	rateHandler := handler.NewRateHandler(c, rateService)
	emailHandler := handler.NewEmailHandler(c, emailService)

	s.Router.Get("/rate", rateHandler.GetCurrentRateHandler())
	s.Router.Post("/subscribe", emailHandler.SubscribeEmailHandler())
	s.Router.Post("/sendEmails", emailHandler.SendToEmailsHandler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port), s.Router); err != nil {
		log.Fatalf(fmt.Sprintf("Failed to listen and serve server on the port - %d", c.Server.Port), err)
	}
}

func setupRateProvider(c *config.Config, logger service.Logger) *provider.CachedCryptoProvider {
	chainOfProviders := provider.SetupChainOfProviders(c.Crypto, logger)
	cacheDuration, err := time.ParseDuration(c.Crypto.CacheDuration)
	if err != nil {
		log.Fatalf("Incorrect cache duration! Please check your configuration.")
	}
	cachedChainOfProviders := provider.NewCachedCryptoProvider(cacheDuration, chainOfProviders, provider.NewDefaultTimeProvider())
	return cachedChainOfProviders
}

func setupLogger() service.Logger {
	loggerType := os.Getenv("BROKER")
	switch loggerType {
	case "rabbit":
		return clogger.NewRabbitMqLogger()
	case "kafka":
		return clogger.NewKafkaLogger()
	default:
		log.Fatalf("Usupported type of message-broker like %s", loggerType)
	}
	return nil
}

func NewServer(conf *config.Config) *Server {
	return &Server{
		Config: conf,
		Router: chi.NewRouter(),
	}
}
