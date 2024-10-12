package kafka

type ProducerFactory interface {
	NewProducer(config map[string]any) (Producer, error)
}

type Producer interface {
	Produce(topic string, message Message) error
}

type Message struct {
	Key     string
	Headers map[string][]byte
	Value   []byte
}
