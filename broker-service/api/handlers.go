package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	resty "resty.dev/v3"
)

type BrokerHandler struct {
	restyClient *resty.Client
}

func NewBrokerHandler(client *resty.Client) *BrokerHandler {
	return &BrokerHandler{
		restyClient: client,
	}
}

// What message Broker will send to FE
type BrokerResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// What message Broker will receive from FE
type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// TODO: make FE send this payload to Broker
type LogPayload struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

func (br *BrokerHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (br *BrokerHandler) HitBroker(c *fiber.Ctx) error {
	response := BrokerResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (br *BrokerHandler) HandleSubmission(c *fiber.Ctx) error {
	var requestPayload RequestPayload

	if err := c.BodyParser(&requestPayload); err != nil {
		log.Printf("Error parsing request payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(BrokerResponse{
			Error: true,
			Message: "Invalid request payload",
		})
	}

	switch requestPayload.Action {
	case "auth":
		return makeAuthServiceCall(c, br.restyClient, &requestPayload)
	case "log":
		return c.Status(fiber.StatusOK).JSON(BrokerResponse{
			Error:   true,
			Message: "Log service not implemented yet",
		});

	default:
		return c.Status(fiber.StatusBadRequest).JSON(BrokerResponse{
			Error: true,
			Message: "Cant recognize requested action",
		})
	}
}

func makeAuthServiceCall(c *fiber.Ctx, client *resty.Client, p *RequestPayload) error {
	// Send post request to auth service
	resp, err := client.R().SetBody(p.Auth).SetResult(&BrokerResponse{}).Post("http://auth-service:9090/login")
	if err != nil {
		log.Printf("Error calling auth service: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(BrokerResponse{
            Error:   true,
            Message: "Authentication service not responding",
        })
	}

	if resp.StatusCode() == fiber.StatusUnauthorized {
		return c.Status(fiber.StatusInternalServerError).JSON(BrokerResponse{
			Error:   true,
			Message: "User is not authentificatied",
		})
	} else if resp.StatusCode() != fiber.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(BrokerResponse{
			Error:   true,
			Message: "Error calling auth server",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp.Result().(*BrokerResponse))
}