package handlers

import (
	"example/internal/api/util/response"
	"example/internal/api/util/upload"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/cr_user"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"path/filepath"
)

// HandleLogin authenticate user require username / email and password
func HandleLogin(service cr_user.Service) fiber.Handler {
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

// HandleLogout invalidate login token
func HandleLogout(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

// HandleMe read user login info
func HandleMe(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		me, err := service.FetchProfile(user)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(me))
	}
}

// HandleUpdateMe update user login info
func HandleUpdateMe(service cr_user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var req struct {
			Payload         entities.CrUser `json:"payload"`
			Password        string          `json:"password"`
			ConfirmPassword string          `json:"confirm_password"`
		}
		if err := c.BodyParser(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(response.ErrorResponse(err))
		}

		if req.Password != req.ConfirmPassword {
			c.Status(http.StatusBadRequest)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "password not match with confirm password"})
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
			req.Payload.Avatar = fileNewPath
		}

		item, err := service.UpdateProfile(user, req.Password, &req.Payload)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessResponse(item))
	}
}
