package kafkaadptr

import (
	"errors"
	"log/slog"

	appkafka "github.com/RuanScherer/journey-track-api/application/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerFactory struct{}

func NewProducerFactory() *ProducerFactory {
	return &ProducerFactory{}
}

func (cf *ProducerFactory) NewProducer(config map[string]any) (appkafka.Producer, error) {
	producer := &producer{}
	if config == nil {
		return nil, errors.New("config is required")
	}

	cfg := kafka.ConfigMap{}
	for k, v := range config {
		cfg[k] = v
	}

	kafkaProducer, err := kafka.NewProducer(&cfg)
	if err != nil {
		slog.Error("Error creating kafka producer.", "error", err)
		return nil, err
	}

	producer.producer = kafkaProducer
	return producer, nil
}

type producer struct {
	producer *kafka.Producer
}

func (p *producer) Produce(topic string, message appkafka.Message) error {
	headers := make([]kafka.Header, 0)
	for key, value := range message.Headers {
		headers = append(headers, kafka.Header{
			Key:   key,
			Value: value,
		})
	}

	msg := &kafka.Message{
		Key: []byte(message.Key),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:   message.Value,
		Headers: headers,
	}
	if err := p.producer.Produce(msg, nil); err != nil {
		slog.Error("Error producing kafka message", "error", err)
		return err
	}

	p.producer.Flush(1000)
	return nil
}
