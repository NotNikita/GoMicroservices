package types

// Payload represents the structure of messages received from RabbitMQ
type Payload struct {
    Name string `json:"name"`
    Data string `json:"data"`
}

// Handler defines the interface for processing event payloads
type Handler interface {
    Handle(payload Payload) error
}