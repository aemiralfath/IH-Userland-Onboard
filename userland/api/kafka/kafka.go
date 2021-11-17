package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka interface {
	SendMessage(topic string, message []byte) error
}

type KafkaConfig struct {
	Host   string
	Port   string
	Group  string
	Offset string
}

type KafkaConfluentinc struct {
	Config *kafka.ConfigMap
}

func NewKafka(config KafkaConfig) Kafka {
	return &KafkaConfluentinc{
		Config: &kafka.ConfigMap{
			"bootstrap.servers": fmt.Sprintf("%s:%s", config.Host, config.Port),
		},
	}
}

func (k *KafkaConfluentinc) SendMessage(topic string, message []byte) error {
	producer, err := kafka.NewProducer(k.Config)
	if err != nil {
		return err
	}

	defer producer.Close()

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

	go deliver(producer)

	// Wait for message deliveries before shutting down
	producer.Flush(15 * 1000)

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
