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
	broker := api.NewBrokerHandler(restyClient)
	app := fiber.New()
	api.Routes(app, broker)

	log.Fatal(app.Listen(webPort))
}
