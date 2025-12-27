package service

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"listener-service/registry"
	"listener-service/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	logServiceUrl = "http://logger-service:7070/log"
	rabbitMqHost = "rabbitmq" // rabbitmq / localhost
)

type RabbitMQService struct {
	connection *amqp.Connection
    registry   *registry.HandlerRegistry
    queueName string
}


func NewRabbitMQService(registry *registry.HandlerRegistry) *RabbitMQService {
    // Create a single connection
    // TODO: rabbitmq for prod, localhost for dev
    conn, err := amqp.Dial("amqp://guest:guest@" + rabbitMqHost + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
    log.Println("Connected to RabbitMQ")
    
    service := &RabbitMQService{connection: conn, registry: registry}
    
    // Set up all the queues once during initialization
    service.setupExchange()
    
    return service
}

func (s *RabbitMQService) ShutDown() {
	if err := s.connection.Close(); err != nil {
		log.Fatalf("Failed to close RabbitMQ connection: %s", err)
	}
}

/*
Exchange = Mail Sorting Facility
Exchange: The sorting center that decides which mailboxes get which letters
It doesn't store mail; it just routes it to the right mailboxes
Different types of exchanges have different routing rules
*/
func (s *RabbitMQService) setupExchange() {
    // Create a temporary channel just for setup
    ch, err := s.connection.Channel()
	failOnError(err, "Failed to open a channel")
    defer ch.Close()
    
    // Declare all the queues your application needs
    ch.ExchangeDeclare(
        "events_exchange", // purpose of exchange
        "topic", // Topics are the categories/routing keys used within that facility
        true, // durable
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil, // arguments
    )
	log.Println("Exchange events_exchange declared with <topic> type")

}

/*
Queue = Mailbox
Queue: A mailbox where messages wait to be picked up
Each message sits in the mailbox until someone comes to collect it
Once collected, the message is gone (consumed)

For each call of Listen, a new queue is created
The queue is bound to the exchange with a specific routing key
The queue will be deleted when the connection is closed
*/
func (s *RabbitMQService) Listen(topics []string) error {
    // Get a new channel for this operation
    wg := sync.WaitGroup{}
    ch, err := s.connection.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()
    
    q, err := declareRandomQueue(ch)
    failOnError(err, "Failed to declare a queue")
	log.Printf("Queue %s declared within exchange events_exchange", q.Name)
    
    /* QueueBind binds an exchange to a queue so that publishings
     to the exchange will be routed to the queue when the publishing 
     routing key matches the binding routing key. 
     
     topics are the routing keys that will be used to bind the queue
     Example: log.* will match log.info, log.error, etc.
     */
    for _, topic := range topics {
        err := ch.QueueBind(
            q.Name,
            topic,  // Ключ маршрутизации
            "events_exchange", // exchange
            false,
            nil, // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to bind queue: %w", err)
        }
        log.Printf("Event with %s key will go to %s queue", topic, q.Name)
    }

    messages, err := ch.Consume(
        q.Name, // queue
        "",
        true,
        false,
        false,
        false,
        nil,
    )

    go func(){
        wg.Add(1)
        defer wg.Done()

        for {
            select {
            case <-ch.NotifyClose(make(chan *amqp.Error)):
                log.Println("Channel closed")
                wg.Done()
                return
            case <-ch.NotifyCancel(make(chan string)):
                log.Println("Channel cancelled")
                wg.Done()
                return
            case newMessage := <-messages:
                // Process the message
                log.Printf("Received a message: %s, %s", newMessage.RoutingKey, newMessage.Body)
                var payload types.Payload
                err := json.Unmarshal(newMessage.Body, &payload)
                failOnError(err, "Failed to unmarshal message")

                // Use the registry to process the event
                err = s.registry.ProcessEvent(newMessage.RoutingKey, payload); 
                failOnError(err, "Failed to process event")

            // Add a default case to prevent blocking if no messages
            default:
                // Sleep a tiny bit to prevent CPU spinning
                time.Sleep(10 * time.Millisecond)
        }
        }
    }()

    wg.Wait()
    return nil
}


func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
    return ch.QueueDeclare(
        "",    // name
        false, // durable
        false, // delete when unused
        true, // exclusive
        false, // no-wait
        nil,   // arguments
    )
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}