package logger

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
)

type KafkaConsumer struct{}

func NewKafkaConsumer() *KafkaConsumer {
	return &KafkaConsumer{}
}

func (r *KafkaConsumer) Consume() {
	consumer := r.setupKafkaConsumer()
	defer consumer.Close()

	partitionConsumer := r.setupPartitionConsumer(consumer)
	defer partitionConsumer.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for message := range partitionConsumer.Messages() {
			if strings.HasPrefix(string(message.Value), "[Error]") {
				log.Printf("%s", string(message.Value))
			}
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	wg.Wait()
}

func (r *KafkaConsumer) setupKafkaConsumer() sarama.Consumer {
	brokerURL := os.Getenv("KAFKA_SERVER_URL")
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{brokerURL}, config)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}

func (r *KafkaConsumer) setupPartitionConsumer(consumer sarama.Consumer) sarama.PartitionConsumer {
	topic := os.Getenv("KAFKA_LOG_TOPIC")
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}
	return partitionConsumer
}
