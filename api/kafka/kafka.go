package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka interface {
	NewConsumer() (*kafka.Consumer, error)
	NewProducer() (*kafka.Producer, error)
	Produce(producer *kafka.Producer, topic string, message []byte) error
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
			"group.id":          config.Group,
			"auto.offset.reset": config.Offset,
		},
	}
}

func (k *KafkaConfluentinc) NewConsumer() (*kafka.Consumer, error) {
	return kafka.NewConsumer(k.Config)
}

func (k *KafkaConfluentinc) NewProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(k.Config)
}

func (k *KafkaConfluentinc) Produce(producer *kafka.Producer, topic string, message []byte) error {
	err := producer.Produce(&kafka.Message{
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
