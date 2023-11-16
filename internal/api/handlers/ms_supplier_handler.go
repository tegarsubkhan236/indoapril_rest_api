package handlers

import (
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/ms_supplier"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetSuppliers(service ms_supplier.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			Page  int    `query:"page" default:"1"`
			Limit int    `query:"limit" default:"10"`
			IDs   []uint `query:"ids"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		data, count, err := service.FetchAllSupplier(q.Page, q.Limit, q.IDs)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(data, count))
	}
}

func GetSupplier(service ms_supplier.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.FetchDetailSupplier(id)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func CreateSupplier(service ms_supplier.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req []entities.MsSupplier
		if err := c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.InsertSupplier(&req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func UpdateSupplier(service ms_supplier.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var req entities.MsSupplier
		if err = c.BodyParser(&req); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.UpdateSupplier(id, &req)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func DeleteSupplier(service ms_supplier.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			ID []uint `query:"id"`
		}
		if err := c.QueryParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err := service.DeleteSupplier(req.ID); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}
