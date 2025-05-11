package main

import (
	"broker/api"
	"log"

	"github.com/gofiber/fiber/v2"
	resty "resty.dev/v3"
)

const (
	webPort = ":8080"
)

func main() {
	log.Println("Broker service started")
	restyClient := resty.New()
	defer restyClient.Close()

	    // Init RabbitMQ service with handler registry
    mqService := api.NewRabbitMQService()
    defer mqService.ShutDown()

	broker := api.NewBrokerHandler(restyClient, mqService)
	app := fiber.New()
	api.Routes(app, broker)

	log.Fatal(app.Listen(webPort))
}
