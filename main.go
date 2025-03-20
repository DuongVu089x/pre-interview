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
	// var arr1 [5]int
	// arr1[0] = 1
	// arr1[1] = 2
	// arr1[2] = 3
	// arr1[3] = 4
	// arr1[4] = 5

	// a := new(*int)
	// fmt.Println(a)

	// return

	// lfu.TestLFU()
	// return
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
			Key: "OD0001",
			Value: domain.MessageValue{
				Meta: &domain.MetaData{
					MessageID: "OD0001",
					ServiceID: "orders-service",
					Timestamp: time.Now().UnixNano(),
				},
				MessageCode: "OD0001",
				Payload: map[string]interface{}{
					"order_id": "OD0001",
					"amount":   100.0,
					"status":   "pending",
				},
			},
			Topic: "orders-topic",
		},
		{
			Key: "PA0001",
			Value: domain.MessageValue{
				Meta: &domain.MetaData{
					MessageID: "PA0001",
					ServiceID: "payments-service",
					Timestamp: time.Now().UnixNano(),
				},
				MessageCode: "PA0001",
				Payload: map[string]interface{}{
					"payment_id": "PA0001",
					"amount":     100.0,
					"status":     "pending",
				},
			},
			Topic: "payments-topic",
		},
		{
			Key: "US0001",
			Value: domain.MessageValue{
				Meta: &domain.MetaData{
					MessageID: "US0001",
					ServiceID: "users-service",
					Timestamp: time.Now().UnixNano(),
				},
				MessageCode: "US0001",
				Payload: map[string]interface{}{
					"user_id": "US0001",
					"name":    "John Doe",
					"email":   "john.doe@example.com",
				},
			},
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

	retryConfig := kafka.RetryConfig{
		RetryTopicSuffix:    "-retry",
		DLQTopicSuffix:      "-dlq",
		MaxRetryAttempts:    3,
		RetryBackoffInitial: 1 * time.Second,
		RetryBackoffMax:     1 * time.Minute,
		RetryBackoffFactor:  2,
	}

	consumer, err := kafka.NewRetryableConsumer(consumerConfig, config, retryConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	{
		// Register different handlers for different topics
		err = consumer.RegisterHandler("orders-topic", func(msg domain.Message) error {
			fmt.Printf("Processing order: key: %s, topic: %s, partition: %d, offset: %d\n", msg.Key, msg.Topic, msg.Partition, msg.Offset)
			// Order-specific business logic

			// Simulate random failures (50% chance)
			if time.Now().UnixNano()%2 == 0 {
				return fmt.Errorf("simulated processing error")
			}

			return nil
		})
		if err != nil {
			log.Fatalf("Failed to register handler: %s", err)
		}

		err = consumer.RegisterHandler("payments-topic", func(msg domain.Message) error {
			fmt.Printf("Processing payment: key: %s, topic: %s, partition: %d, offset: %d\n", msg.Key, msg.Topic, msg.Partition, msg.Offset)
			// Payment-specific business logic
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to register handler: %s", err)
		}

		err = consumer.RegisterHandler("users-topic", func(msg domain.Message) error {
			fmt.Printf("Processing user event: key: %s, topic: %s, partition: %d, offset: %d\n", msg.Key, msg.Topic, msg.Partition, msg.Offset)
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

	for i := 1; i <= 10; i++ {
		msg := domain.Message{
			Key: fmt.Sprintf("order-%d", i),
			Value: domain.MessageValue{
				Meta: &domain.MetaData{
					MessageID: fmt.Sprintf("order-%d", i),
					ServiceID: "orders-service",
					Timestamp: time.Now().UnixNano(),
				},
				MessageCode: fmt.Sprintf("OD000%d", i),
				Payload: map[string]interface{}{
					"order_id": fmt.Sprintf("OD000%d", i),
					"amount":   100.0,
					"status":   "pending",
				},
			},
			Topic: "orders-topic",
		}
		producer.Publish(msg)
	}

	// Wait for termination
	<-ctx.Done()
	fmt.Println("Consumer shutdown complete")
}
