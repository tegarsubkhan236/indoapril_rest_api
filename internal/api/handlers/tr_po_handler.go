package handlers

import (
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/tr_purchase_order"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func HandleGetPurchaseOrders(service tr_purchase_order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			Page  int `query:"page" default:"1"`
			Limit int `query:"limit" default:"10"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		data, count, err := service.FetchAll(q.Page, q.Limit)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(data, count))
	}
}

func HandleGetPurchaseOrder(service tr_purchase_order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			PoCode *string `query:"po_code"`
			BoCode *string `query:"bo_code"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		itemPO, itemBO, err := service.FetchDetail(q.PoCode, q.BoCode)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		if itemPO != nil {
			return c.JSON(response.SuccessResponse(itemPO))
		} else {
			return c.JSON(response.SuccessResponse(itemBO))
		}
	}
}

func HandleAddPurchaseOrder(service tr_purchase_order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.TrPurchaseOrderReq
		if err := c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.Insert(req)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}
