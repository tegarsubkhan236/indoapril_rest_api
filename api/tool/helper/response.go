package helper

import "github.com/gofiber/fiber/v2"

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseHandler(c *fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(response{
		Status:  code,
		Message: message,
		Data:    data,
	})
}
