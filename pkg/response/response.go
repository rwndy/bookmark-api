package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func JSON(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{Data: data, Message: message, Status: status})
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusOK, message, data)
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusCreated, message, data)
}

func Fail(c *fiber.Ctx, status int, message string) error {
	return JSON(c, status, message, nil)
}
