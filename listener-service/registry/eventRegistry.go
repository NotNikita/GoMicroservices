package registry

import (
	"fmt"
	"log"

	"listener-service/types" // Updated import
)

// HandlerRegistry maintains a mapping of event types to their handlers
type HandlerRegistry struct {
	handlers map[string]types.Handler // Updated type
}

// NewHandlerRegistry creates a new HandlerRegistry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]types.Handler),
	}
}

// Register adds a handler for a specific event type
func (r *HandlerRegistry) Register(eventType string, handler types.Handler) {
	r.handlers[eventType] = handler
}

// GetHandler retrieves the appropriate handler for an event type
func (r *HandlerRegistry) GetHandler(eventType string) types.Handler {
	handler, exists := r.handlers[eventType]
	if !exists {
		log.Printf("No handler registered for event type: %s", eventType)
		return nil
	}
	return handler
}

// ProcessEvent processes an event payload using the appropriate handler
func (r *HandlerRegistry) ProcessEvent(routingKey string, payload types.Payload) error {
	handler := r.GetHandler(routingKey)
	if handler == nil {
		return fmt.Errorf("no handler for event type: %s, payload: %s", routingKey, payload)
	}

	return handler.Handle(payload)
}