package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Routes(app *fiber.App, broker *BrokerHandler) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Link",
		// AllowCredentials: true,
		MaxAge: 300,
	}))

	app.Get("/ping", broker.HealthCheck)
	app.Post("/", broker.HitBroker)
	app.Post("/handle", broker.HandleSubmission)
}
