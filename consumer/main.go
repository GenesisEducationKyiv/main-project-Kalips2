package main

import (
	domain "consumer/domain"
	"consumer/logger"
	"log"
	"os"
)

func main() {
	consumer := createConsumer()
	consumer.Consume()
}

func createConsumer() domain.Consumer {
	loggerType := os.Getenv("BROKER")
	switch loggerType {
	case "rabbit":
		return logger.NewRabbitConsumer()
	case "kafka":
		return logger.NewKafkaConsumer()
	default:
		log.Fatalf("Usupported type of message-broker like %s", loggerType)
	}
	return nil
}
