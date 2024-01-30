package handlers

import (
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/ms_product"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func HandleGetProducts(service ms_product.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			Page              int    `query:"page" default:"1"`
			Limit             int    `query:"limit" default:"10"`
			ProductName       string `query:"product_name"`
			SupplierID        uint   `query:"supplier_id"`
			ProductCategories []uint `query:"product_categories"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		data, count, err := service.FetchAllProduct(q.Page, q.Limit, q.ProductName, q.SupplierID, q.ProductCategories)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(data, count))
	}
}

func HandleGetProduct(service ms_product.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var basedOn = c.Params("based_on") // Price or Stock

		item, err := service.FetchDetailProduct(id, basedOn)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleAddProduct(service ms_product.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req []entities.MsProductReq
		if err := c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		for _, item := range req {
			if err := item.ValidateInput(); err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(response.ErrorResponse(err))
			}
		}

		data, err := service.InsertBatchProduct(&req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(data))
	}
}

func HandleUpdateProduct(service ms_product.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var payload entities.MsProductReq
		if err = c.BodyParser(&payload); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = payload.ValidateInput(); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.UpdateProduct(id, &payload)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleRemoveProduct(service ms_product.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			ID []uint `query:"id"`
		}
		if err := c.QueryParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		err := service.DeleteProduct(req.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}
