package api

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	logServiceUrl = "http://logger-service:7070/log"
	rabbitMqHost = "rabbitmq" // rabbitmq / localhost
	exchangeName = "events_exchange"
)

type RabbitMQService struct {
	connection *amqp.Connection
}


func NewRabbitMQService() *RabbitMQService {
    // Create a single connection
    conn, err := amqp.Dial("amqp://guest:guest@" + rabbitMqHost + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
    log.Println("Connected to RabbitMQ")
    
    service := &RabbitMQService{connection: conn}
    
    return service
}

func (s *RabbitMQService) ShutDown() {
	if err := s.connection.Close(); err != nil {
		log.Fatalf("Failed to close RabbitMQ connection: %s", err)
	}
}

func (s *RabbitMQService) EmitRabbitMQMessage(c *fiber.Ctx, p *RequestPayload) error {
	// Create channel
	ch, err := s.connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	failOnError(err, "Failed to re-declare exchange")

	jsonData, err := json.Marshal(p.Log)
	failOnError(err, "Failed to marshal LogPayload")

	log.Println("Channel opened. Pushing message to RabbitMQ")
	err = ch.Publish(exchangeName, "log.INFO", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
	})
	failOnError(err, "Failed to publish a message")

	log.Println("Message sent to RabbitMQ")
	return c.Status(fiber.StatusOK).JSON(BrokerResponse{
		Error:   false,
		Message: "Message sent to RabbitMQ",
	})
}

func failOnError(err error, msg string)error  {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		return err
	}
	return nil
}