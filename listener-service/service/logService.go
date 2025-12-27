package service

import (
	"fmt"
	"listener-service/types"
	"log"

	"resty.dev/v3"
)

// Handler defines the interface for processing event payloads
type Handler interface {
	Handle(payload types.Payload) error
}

const logServiceURL = "http://logger-service:7070/log" // logger-service / localhost

type EventService struct {
	restyClient *resty.Client
}

// NewEventService creates a new EventService
func NewEventService(client *resty.Client) *EventService {
	return &EventService{
		restyClient: client,
	}
}

// Handle processes log and event payloads
func (es *EventService) Handle(payload types.Payload) error {
	type LogResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp, err := es.restyClient.R().
		SetBody(payload).
		SetHeader("Content-Type", "application/json").
		SetResult(&LogResponse{}).
		Post(logServiceURL)

	if err != nil {
		log.Printf("Error calling logger service: %v", err)
		return err
	}

	log.Printf("Logger response status: %d", resp.StatusCode())

	if resp.StatusCode() != 200 {
		result := resp.Result().(*LogResponse)
		log.Printf("Logger service error: %s", result.Message)
		return fmt.Errorf("failed to log request: %s", result.Message)
	}

	return nil
}
