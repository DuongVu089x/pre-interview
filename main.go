package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DuongVu089x/pre-interview/domain"
	"github.com/DuongVu089x/pre-interview/infrastructure/kafka"
)

// func subarraysWithKDistinct(nums []int, k int) int {

// 	windowMap := make(map[int]map[int]interface{})
// 	result := 0
// 	// currentWindow := []int

// 	for right := 0; right < len(nums)-1; right++ {

// 		for left := right + 1; left < len(nums); left++ {
// 			if numsA, ok := windowMap[right]; ok {
// 				numsA[nums[left]] = struct{}{}
// 				if len(numsA) >= k {
// 					break
// 				}
// 			} else {
// 				windowMap[right] = make(map[int]interface{})
// 				windowMap[right][nums[left]] = struct{}{}
// 			}

// 			result++
// 		}
// 	}

// 	return result
// }

func main() {
	// fmt.Println(lengthOfLongestSubstring("abcbacbb"))

	// Initialize Kafka producer with config from your YAML
	config := kafka.ProducerConfig{
		BootstrapServers: "localhost:9092",
		SecurityProtocol: "plaintext",
		DefaultTopic:     "first-topic",
		Topics: []kafka.TopicConfig{
			{
				Name:              "orders-topic",
				NumPartitions:     3,
				ReplicationFactor: 3,
			},
			{
				Name:              "payments-topic",
				NumPartitions:     3,
				ReplicationFactor: 3,
			},
			{
				Name:              "users-topic",
				NumPartitions:     3,
				ReplicationFactor: 3,
			},
		},
	}

	// Create topics before producing
	err := kafka.CreateTopics(config)
	if err != nil {
		log.Printf("Warning: Topic creation failed: %s", err)
		// Continue anyway - topics might already exist
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Create and publish a message
	messages := []domain.Message{
		{
			Key:   "OD0001",
			Value: "Order 1",
			Topic: "orders-topic",
		},
		{
			Key:   "PA0001",
			Value: "Payment 1",
			Topic: "payments-topic",
		},
		{
			Key:   "US0001",
			Value: "User 1",
			Topic: "users-topic",
		},
	}

	for _, message := range messages {
		err = producer.Publish(message)
		if err != nil {
			log.Fatalf("Failed to publish message: %s", err)
		}
	}

	// Give some time for the message to be delivered
	// In a real application, you might want to implement proper waiting
	time.Sleep(1 * time.Second)

	// Initialize Kafka consumer with config from your YAML
	consumerConfig := kafka.ConsumerConfig{
		BootstrapServers: "localhost:9092",
		SecurityProtocol: "plaintext",
		GroupID:          "my-consumer-group",
		AutoOffsetReset:  "earliest",
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	{
		// Register different handlers for different topics
		err = consumer.RegisterHandler("orders-topic", func(msg domain.Message) error {
			fmt.Printf("Processing order: %s, key: %s, topic: %s, partition: %d, offset: %d\n", msg.Value, msg.Key, msg.Topic, msg.Partition, msg.Offset)
			// Order-specific business logic
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to register handler: %s", err)
		}

		err = consumer.RegisterHandler("payments-topic", func(msg domain.Message) error {
			fmt.Printf("Processing payment: %s, key: %s, topic: %s, partition: %d, offset: %d\n", msg.Value, msg.Key, msg.Topic, msg.Partition, msg.Offset)
			// Payment-specific business logic
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to register handler: %s", err)
		}

		err = consumer.RegisterHandler("users-topic", func(msg domain.Message) error {
			fmt.Printf("Processing user event: %s, key: %s, topic: %s, partition: %d, offset: %d\n", msg.Value, msg.Key, msg.Topic, msg.Partition, msg.Offset)
			// User-specific business logic
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to register handler: %s", err)
		}
	}

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Received termination signal, shutting down...")
		cancel()
	}()

	// Start consuming messages in a goroutine
	go func() {
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			log.Fatalf("Consumer error: %s", err)
		}
	}()

	// Wait for termination
	<-ctx.Done()
	fmt.Println("Consumer shutdown complete")
}
