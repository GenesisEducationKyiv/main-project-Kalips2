package logger

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RabbitLogger struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (l *RabbitLogger) LogDebug(v ...any) error {
	return l.publishMessage(Debug, fmt.Sprint(v...))
}

func (l *RabbitLogger) LogError(v ...any) error {
	return l.publishMessage(Error, fmt.Sprint(v...))
}

func (l *RabbitLogger) LogInfo(v ...any) error {
	return l.publishMessage(Info, fmt.Sprint(v...))
}

func (l *RabbitLogger) publishMessage(logType string, message string) error {
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}
	return l.channel.Publish(
		"",      // exchange
		logType, // routing key
		false,   // mandatory
		false,   // immediate
		msg)
}

func NewRabbitMqLogger() *RabbitLogger {
	connection := setupAMQPServer()
	channel := setupAMQPChanel(connection)
	declareQueue(channel, Debug)
	declareQueue(channel, Error)
	declareQueue(channel, Info)
	return &RabbitLogger{connection, channel}
}

func setupAMQPServer() *amqp.Connection {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	connection, err := amqp.Dial(amqpServerURL)
	if err != nil {
		log.Fatal(err)
	}
	return connection
}

func setupAMQPChanel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return channel
}

func declareQueue(channel *amqp.Channel, name string) {
	_, err := channel.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
}
