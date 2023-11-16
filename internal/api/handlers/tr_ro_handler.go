package handlers

import (
	"errors"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/tr_receiving_order"
	"github.com/gofiber/fiber/v2"
)

func CreateReceivingOrder(service tr_receiving_order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.TrReceivingOrderReq
		if err := c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if req.PoCode == "" {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(errors.New("po code cannot be empty")))
		}

		if req.UserID == 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(errors.New("user id cannot be empty")))
		}

		item, err := service.Create(req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(fiber.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}
