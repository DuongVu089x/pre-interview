package kafka

import (
	"fmt"

	"github.com/DuongVu089x/pre-interview/application/port"
	"github.com/DuongVu089x/pre-interview/domain"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer implements the MessageProducer interface using Kafka
type Producer struct {
	producer     *kafka.Producer
	defaultTopic string
}

// NewProducer creates a new Kafka producer
func NewProducer(config ProducerConfig) (*Producer, error) {
	p, err := kafka.NewProducer(config.NewConfigMap())
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Start a goroutine to handle delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()

	return &Producer{
		producer:     p,
		defaultTopic: config.DefaultTopic,
	}, nil
}

// Publish sends a message to Kafka
func (p *Producer) Publish(message domain.Message) error {
	topic := message.Topic
	if topic == "" {
		topic = p.defaultTopic
	}

	// Create delivery channel for this specific message
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(message.Partition),
		},
		Key:   []byte(message.Key),
		Value: []byte(message.Value),
	}

	// Produce with delivery channel
	if err := p.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("message delivery failed: %w", m.TopicPartition.Error)
	}

	return nil
}

// Close flushes and closes the producer
func (p *Producer) Close() error {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
	return nil
}

// Ensure Producer implements MessageProducer interface
var _ port.MessageProducer = (*Producer)(nil)
