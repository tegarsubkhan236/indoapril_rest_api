package request

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func IdParam(c *fiber.Ctx) (ID uint, err error) {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
