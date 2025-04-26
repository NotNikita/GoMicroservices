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

func (lh *LoggerHandler) UpdateLog(c *fiber.Ctx) error {
	var logEvent service.LogEntry

	if err := c.BodyParser(&logEvent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	res, err := lh.ls.UpdateOne(logEvent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update log entry",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Log entry was updated successfully",
		"data":    res,
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Logs fetched successfully",
		"data":    logs,
	})
}

func (lh *LoggerHandler) GetLogById(c *fiber.Ctx) error {
	id := c.Params("id")
	log, err := lh.ls.GetOne(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch log",
		})
	}
	if log == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Log not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Log fetched successfully",
		"data":    log,
	})
}

func (lh *LoggerHandler) ClearLogs(c *fiber.Ctx) error {
	err := lh.ls.DropCollection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Logs cleared successfully",
	})
}