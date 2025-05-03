package api

import (
	"auth/model"
	"auth/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	resty "resty.dev/v3"
)

const (
	logServiceUrl = "http://logger-service:7070/log"
)

type UserHandler struct {
	service *service.UserService
	restyClient *resty.Client
}

func NewUserHandler(service *service.UserService, restyClient *resty.Client) *UserHandler {
	return &UserHandler{
		service: service,
		restyClient: restyClient,
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
		log.Printf("post.Login password hash didnt match for: %s %v %v", userPayload.Email, err, result)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email of password is incorrect",
		})
	}
	
	err = uh.logRequest("authentication", fmt.Sprintf("User %s logged in", userPayload.Email))
	if err != nil {
		log.Printf("post.Login failed to log successfull auth call: %v", err)
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": fmt.Sprintf("Logged in user %s", userPayload.Email),
		"data":    user,
	})
}

func (uh *UserHandler) logRequest(name, data string) error {
    payload := struct {
        Name string `json:"name"`
        Data string `json:"data"`
    }{
        Name: name,
        Data: data,
    }

    type LogResponse struct {
        Error   bool   `json:"error"`
        Message string `json:"message"`
    }
    

    resp, err := uh.restyClient.R().
        SetBody(payload).
        SetHeader("Content-Type", "application/json").
        SetResult(&LogResponse{}).
        Post(logServiceUrl)
        
    if err != nil {
        log.Printf("Error calling logger service: %v", err)
        return err
    }

    log.Printf("Logger response status: %d", resp.StatusCode())
    
    if resp.StatusCode() != fiber.StatusOK {
        result := resp.Result().(*LogResponse)
        log.Printf("Logger service error: %s", result.Message)
        return fmt.Errorf("failed to log request: %s", result.Message)
    }

    return nil
}