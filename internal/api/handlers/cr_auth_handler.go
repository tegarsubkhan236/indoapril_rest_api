package handlers

import (
	"errors"
	"example/internal/api/util/response"
	"example/internal/api/util/upload"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/cr_user"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"path/filepath"
)

// Login authenticate user require username / email and password
func Login(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginInput struct {
			Identity string `json:"identity"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&loginInput); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		token, err := service.AuthenticateUser(loginInput.Identity, loginInput.Password)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(token))
	}
}

// Me read user login info
func Me(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		me, err := service.FetchProfile()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(me))
	}
}

// UpdateMe update user login info
func UpdateMe(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.CrUser
		if err := c.BodyParser(&req); err != nil {
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

		item, err := service.UpdateProfile(&req)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}

// UpdatePasswordMe update user password login info
func UpdatePasswordMe(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			OldPassword     string `json:"old_password"`
			NewPassword     string `json:"new_password"`
			ConfirmPassword string `json:"confirm_password"`
		}
		if err := c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if req.NewPassword != req.ConfirmPassword {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(errors.New("password confirmation does not match")))
		}

		item, err := service.UpdateProfilePassword(req.OldPassword, req.NewPassword)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}