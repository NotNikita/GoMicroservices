package main

import (
	"broker/api"
	"log"

	"github.com/gofiber/fiber/v2"
)

const (
	webPort = ":8080"
)

func main() {
	log.Println("Broker service started")
	app := fiber.New()
	api.Routes(app)

	app.Get("/ping", api.HealthCheck)
	app.Post("/", api.Broker)

	log.Fatal(app.Listen(webPort))
}
