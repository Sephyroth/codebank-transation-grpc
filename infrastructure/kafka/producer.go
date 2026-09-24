package kafka

import (
	"fmt"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer *ckafka.Producer
}

func NewKafkaProducer() KafkaProducer {
	return KafkaProducer{}
}

func (k *KafkaProducer) SetupProducer(bootstrapServer string) error {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
	}
	for _, key := range []string{"security.protocol", "sasl.mechanisms", "sasl.username", "sasl.password"} {
		if value := os.Getenv(key); value != "" {
			if err := configMap.SetKey(key, value); err != nil {
				return err
			}
		}
	}
	producer, err := ckafka.NewProducer(configMap)
	if err != nil {
		return err
	}
	k.Producer = producer
	return nil
}

func (k *KafkaProducer) Publish(msg string, topic string) error {
	if k.Producer == nil {
		return fmt.Errorf("Kafka producer is not initialized")
	}
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          []byte(msg),
	}
	err := k.Producer.Produce(message, nil)
	if err != nil {
		return err
	}
	return nil
}
