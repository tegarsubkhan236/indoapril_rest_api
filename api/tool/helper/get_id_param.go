package helper

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetIDParam(c *fiber.Ctx) (uint, error) {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
