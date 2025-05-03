package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Routes(app *fiber.App, logger *LoggerHandler) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Link",
		// AllowCredentials: true,
		MaxAge: 300,
	}))

	app.Get("/ping", logger.HealthCheck)
	app.Post("/log", logger.CreateLog)
	app.Put("/log", logger.UpdateLog)
	app.Get("/logs", logger.GetAllLogs)
	app.Get("/log/:id", logger.GetLogById)
	app.Delete("/logs/drop", logger.ClearLogs)
}
