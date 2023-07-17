package logger

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RabbitConsumer struct{}

func NewRabbitConsumer() *RabbitConsumer {
	return &RabbitConsumer{}
}

func (r *RabbitConsumer) Consume() {
	connection := r.setUpAMQPConnection()
	defer connection.Close()

	channel := r.setUpAMQPChannel(connection)
	defer channel.Close()

	queue := os.Getenv("CONSUME_LOG_LEVEL")
	messages := r.setUpMessage(channel, queue)

	forever := make(chan bool)
	go func() {
		for message := range messages {
			log.Printf("[%s]: %s", queue, message.Body)
		}
	}()

	<-forever
}

func (r *RabbitConsumer) setUpMessage(channel *amqp.Channel, queue string) <-chan amqp.Delivery {
	message, err := channel.Consume(
		queue, // queue name
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	return message
}

func (r *RabbitConsumer) setUpAMQPConnection() *amqp.Connection {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	connection, err := amqp.Dial(amqpServerURL)
	if err != nil {
		log.Fatal(err)
	}
	return connection
}

func (r *RabbitConsumer) setUpAMQPChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return channel
}
