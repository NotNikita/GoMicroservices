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
	app := fiber.New()
	api.Routes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/ping", api.HealthCheck)
	app.Post("/", api.Broker)

	log.Fatal(app.Listen(webPort))
}
