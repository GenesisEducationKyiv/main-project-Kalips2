package logger

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

type KafkaLogger struct {
	writer *kafka.Writer
}

func (l *KafkaLogger) publishMessage(logType string, message string) {
	err := l.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(logType),
			Value: []byte(message),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func NewKafkaLogger() *KafkaLogger {
	kafkaBrokerURL := os.Getenv("KAFKA_BROKER_URL")
	writer := &kafka.Writer{
		Addr:  kafka.TCP(kafkaBrokerURL),
		Topic: "logs",
	}

	return &KafkaLogger{writer}
}
