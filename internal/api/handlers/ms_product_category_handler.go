package handlers

import (
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/ms_product_category"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetProductCategories(service ms_product_category.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			Page  int `q:"page" default:"1"`
			Limit int `q:"limit" default:"10"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		data, count, err := service.FetchAllProductCategory(q.Page, q.Limit)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(data, count))
	}
}

func GetProductCategory(service ms_product_category.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.FetchDetailProductCategory(id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func CreateProductCategory(service ms_product_category.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.MsProductCategory
		if err := c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		err := service.InsertProductCategory(&req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}

func UpdateProductCategory(service ms_product_category.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var req entities.MsProductCategory
		if err = c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = service.UpdateProductCategory(id, &req); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}

func DeleteProductCategory(service ms_product_category.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = service.DeleteProductCategory(id); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}
