package api

import (
	"logger/service"

	"github.com/gofiber/fiber/v2"
)

type LoggerHandler struct {
	ls *service.LoggerService
}

func NewLoggerHandler(ls *service.LoggerService) *LoggerHandler {
	return &LoggerHandler{
		ls: ls,
	}
}

func (lh *LoggerHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (lh *LoggerHandler) CreateLog(c *fiber.Ctx) error {
	var logEvent service.LogEntry

	if err := c.BodyParser(&logEvent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := lh.ls.Insert(logEvent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert log entry",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Log entry created successfully",
		"data":    nil,
	})
}

func (lh *LoggerHandler) GetAllLogs(c *fiber.Ctx) error {
	var queryLimit int64

	if limit := c.Query("limit"); limit != "" {
		queryLimit  = int64(c.QueryInt(limit))
	}

	logs, err := lh.ls.GetAll(queryLimit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch logs",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Logs fetched successfully",
		"data":    logs,
	})
}