package api

import (
	"github.com/gofiber/fiber/v2"
)

type BrokerResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func Broker(c *fiber.Ctx) error {
	response := BrokerResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
