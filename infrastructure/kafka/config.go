package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

// TopicConfig holds configuration for a Kafka topic
type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
}

// Config holds Kafka configuration
type ProducerConfig struct {
	BootstrapServers string
	SecurityProtocol string
	DefaultTopic     string
	Topics           []TopicConfig // Add topic configurations
}

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	BootstrapServers string
	SecurityProtocol string
	GroupID          string
	AutoOffsetReset  string
}

// NewConfigMap converts our config to Kafka's ConfigMap
func (c *ProducerConfig) NewConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": c.BootstrapServers,
		"security.protocol": c.SecurityProtocol,
		"acks":              "all",
	}
}

// NewConsumerConfigMap converts our config to Kafka's ConfigMap
func (c *ConsumerConfig) NewConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  c.BootstrapServers,
		"security.protocol":  c.SecurityProtocol,
		"group.id":           c.GroupID,
		"auto.offset.reset":  c.AutoOffsetReset,
		"enable.auto.commit": true,
	}
}
