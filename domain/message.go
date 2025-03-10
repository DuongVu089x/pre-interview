package domain

type Message struct {
	Key   string
	Value string
	Topic string

	// Add fields for consumer if needed:
	Partition int
	Offset    int64
}
