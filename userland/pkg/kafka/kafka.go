package kafka

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka interface {
	SendMessage(topic string, message []byte) error
}

type KafkaConfluentinc struct {
	Config   *kafka.ConfigMap
	Producer *kafka.Producer
}

func NewKafka() (Kafka, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT")),
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &KafkaConfluentinc{
		Config:   config,
		Producer: producer,
	}, nil
}

func (k *KafkaConfluentinc) SendMessage(topic string, message []byte) error {
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

	go deliver(k.Producer)

	// Wait for message deliveries before shutting down
	// defer k.Producer.Close()

	k.Producer.Flush(15 * 1000)

	return err
}

func deliver(producer *kafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}
