package port

import (
	"context"

	"github.com/DuongVu089x/pre-interview/domain"
)

// MessageConsumer defines the interface for receiving messages
type MessageConsumer interface {
	// Subscribe(topics []string) error
	// Consume(ctx context.Context, handler func(domain.Message) error) error

	RegisterHandler(topic string, handler func(domain.Message) error) error
	Start(ctx context.Context) error
	Close() error
}
