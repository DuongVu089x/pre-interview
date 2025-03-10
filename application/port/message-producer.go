package port

import "github.com/DuongVu089x/pre-interview/domain"

// MessageProducer defines the interface for sending messages
type MessageProducer interface {
	Publish(message domain.Message) error
	Close() error
}
