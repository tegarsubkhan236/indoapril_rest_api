package handlers

import (
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/tr_sales_order"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func HandleGetSalesOrders(service tr_sales_order.Service) fiber.Handler {
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

func HandleGetSalesOrder(service tr_sales_order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.FetchDetail(id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleAddSalesOrder(service tr_sales_order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.TrSalesOrderReq
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
