package handlers

import (
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/cr_permission"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func HandleGetPermissions(service cr_permission.Service) fiber.Handler {
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

func HandleGetPermission(service cr_permission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.FetchDetail(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleAddPermission(service cr_permission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.CrPermission
		if err := c.BodyParser(&requestBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err := requestBody.ValidateInput(); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err := service.Insert(&requestBody); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}

func HandleUpdatePermission(service cr_permission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var requestBody entities.CrPermission
		if err = c.BodyParser(&requestBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = requestBody.ValidateInput(); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err = service.Update(id, &requestBody); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}

func HandleRemovePermission(service cr_permission.Service) fiber.Handler {
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
