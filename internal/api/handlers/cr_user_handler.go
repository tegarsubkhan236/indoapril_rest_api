package handlers

import (
	"errors"
	"example/internal/api/util/request"
	"example/internal/api/util/response"
	"example/internal/api/util/upload"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/cr_user"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"path/filepath"
)

func HandleGetUsers(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q struct {
			Page  int `q:"page" default:"1"`
			Limit int `q:"limit" default:"10"`
		}
		if err := c.QueryParser(&q); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		data, count, err := service.FetchAllUser(q.Page, q.Limit)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(data, count))
	}
}

func HandleGetUser(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.FetchDetailUser(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleAddUser(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.CrUser
		if err := c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		item, err := service.InsertUser(&req)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleUpdateUser(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var req entities.CrUser
		if err = c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		file, err := c.FormFile("avatar")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if file != nil {
			fileNewPath := filepath.Join(upload.Dir, file.Filename)
			if err = c.SaveFile(file, fileNewPath); err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(response.ErrorResponse(err))
			}
			req.Avatar = fileNewPath
		}

		item, err := service.UpdateUser(id, &req)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleUpdateUserPassword(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := request.IdParam(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		var req struct {
			OldPassword     string `json:"old_password"`
			NewPassword     string `json:"new_password"`
			ConfirmPassword string `json:"confirm_password"`
		}
		if err = c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if req.NewPassword != req.ConfirmPassword {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(errors.New("password confirmation does not match")))
		}

		item, err := service.UpdateUserPassword(id, req.OldPassword, req.NewPassword)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

func HandleRemoveUser(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var batchQuery struct {
			ID []uint `q:"id"`
		}
		if err := c.QueryParser(&batchQuery); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if err := service.DeleteUser(batchQuery.ID); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(nil))
	}
}
