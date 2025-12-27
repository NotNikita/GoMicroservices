package main

import (
	"listener-service/registry"
	"listener-service/service"

	"log"
	"time"

	"resty.dev/v3"
)

const (
	port = ":6060"
)

func main(){
    // Init resty client
    restyClient := resty.New()
    defer restyClient.Close()
    
    // Set up handler registry
    register := registry.NewHandlerRegistry()
    
    // Register handlers
    logHandler := service.NewEventService(restyClient)
    register.Register("log.INFO", logHandler)
    register.Register("log.ERROR", logHandler)
    register.Register("log.WARNING", logHandler)
    
    // Init RabbitMQ service with handler registry
    mqService := service.NewRabbitMQService(register)
    defer mqService.ShutDown()
    
    // Start listening for events
    topics := []string{"log.*", "event.*"}
    err := mqService.Listen(topics)
    if err != nil {
        log.Fatalf("Failed to start listening: %v", err)
    }
    
    log.Println("Listener service started")
    time.Sleep(5 * time.Hour)
}
