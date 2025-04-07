package api

import (
	"auth/model"
	"auth/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) Health(c *fiber.Ctx) error {
	return c.SendString("Pong")
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var userPayload model.UserPayload

	if err := c.BodyParser(&userPayload); err != nil {
		log.Println("post.Login Failed to parse body of request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// validate user against database
	user, err := uh.service.GetByEmail(userPayload.Email)
	if err != nil {
		log.Printf("post.Login couldn't find user by email: %s\n", userPayload.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Couldn't find user",
		})
	}

	result, err := uh.service.PasswordMatches(user, userPayload.Password)
	if err != nil || !result {
		log.Printf("post.Login password hash didnt match for: %s\n", userPayload.Email)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Email of password is incorrect",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Error":   false,
		"Message": fmt.Sprintf("Logged in user %s", userPayload.Email),
		"Data":    user,
	})
}
