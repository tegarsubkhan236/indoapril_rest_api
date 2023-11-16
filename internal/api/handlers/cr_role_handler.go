package handlers

import (
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/cr_role"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetRoles(service cr_role.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			Page  int `query:"page" default:"1"`
			Limit int `query:"limit" default:"10"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		data, count, err := service.FetchAll(q.Page, q.Limit)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(data, count))
	}
}

func GetRole(service cr_role.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.FetchDetail(id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func AddRole(service cr_role.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.CrRole
		if err := c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err := service.Insert(&req); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}

func UpdateRole(service cr_role.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var req entities.CrRole
		if err = c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = service.Update(id, &req); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}

func RemoveRole(service cr_role.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = service.Delete(id); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}
