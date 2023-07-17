package logger

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

type KafkaLogger struct {
	producer sarama.SyncProducer
	logTopic string
}

func (l *KafkaLogger) LogDebug(v ...any) error {
	return l.publishToLogs("[" + Debug + "]:" + "[" + fmt.Sprint(v...) + "]")
}

func (l *KafkaLogger) LogError(v ...any) error {
	return l.publishToLogs("[" + Error + "]:" + "[" + fmt.Sprint(v...) + "]")
}

func (l *KafkaLogger) LogInfo(v ...any) error {
	return l.publishToLogs("[" + Info + "]:" + "[" + fmt.Sprint(v...) + "]")
}

func (l *KafkaLogger) publishToLogs(message string) error {
	return l.publishMessageToTopic(l.logTopic, message)
}

func (l *KafkaLogger) publishMessageToTopic(topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := l.producer.SendMessage(msg)
	return err
}

func NewKafkaLogger() *KafkaLogger {
	producer := setupKafkaProducer()
	logTopic := os.Getenv("KAFKA_LOG_TOPIC")
	return &KafkaLogger{producer: producer, logTopic: logTopic}
}

func setupKafkaProducer() sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	brokerURL := os.Getenv("KAFKA_SERVER_URL")
	producer, err := sarama.NewSyncProducer([]string{brokerURL}, cfg)
	if err != nil {
		log.Fatal(err)
	}
	return producer
}
